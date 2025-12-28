import React, { useState, useRef, useEffect } from 'react';
import {
  Drawer,
  Box,
  Typography,
  IconButton,
  TextField,
  List,
  ListItem,
  Paper,
  InputAdornment,
  CircularProgress,
  FormControl,
  Select,
  MenuItem,
} from '@mui/material';
import type { SelectChangeEvent } from '@mui/material';
import {
  Close as CloseIcon,
  Send as SendIcon,
  SmartToy as BotIcon,
} from '@mui/icons-material';
import api from '../../api/client';

interface Message {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: number;
}

interface AgentChatSidebarProps {
  open: boolean;
  onClose: () => void;
}

const ROLES = [
  {
    id: 'reader',
    label: 'General Assistant',
    type: 'reader',
    description: 'Default assistant for reading help',
    profile: {
      avatar: '',
      language: 'Chinese',
      bio: 'A helpful reading assistant.'
    }
  },
  {
    id: 'sherlock',
    label: 'Sherlock Holmes',
    type: 'character',
    description: 'The famous detective.',
    profile: {
      avatar: '',
      language: 'English',
      bio: 'I am Sherlock Holmes, the world\'s only consulting detective. I observe and deduce.',
      book_name: 'Sherlock Holmes'
    }
  },
  {
    id: 'historian',
    label: 'History Teacher',
    type: 'historical',
    description: 'An expert in history.',
    profile: {
      avatar: '',
      language: 'Chinese',
      bio: 'I am a history teacher who loves to explain historical context.',
      historical_era: 'General History'
    }
  }
];

const AgentChatSidebar: React.FC<AgentChatSidebarProps> = ({ open, onClose }) => {
  const [messages, setMessages] = useState<Message[]>([
    {
      id: '1',
      role: 'assistant',
      content: '你好！我是你的智能阅读助手。有什么可以帮你的吗？',
      timestamp: Date.now()
    }
  ]);
  const [inputValue, setInputValue] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [agentId, setAgentId] = useState<string | null>(localStorage.getItem('agent_id'));
  const [selectedRole, setSelectedRole] = useState(ROLES[0].id);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const abortControllerRef = useRef<AbortController | null>(null);
  const bufferRef = useRef('');

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  // Ensure an agent exists or update it if role changes
  useEffect(() => {
    const initOrUpdateAgent = async () => {
      const roleConfig = ROLES.find(r => r.id === selectedRole) || ROLES[0];
      
      if (!agentId) {
        // Create new agent
        try {
          const res = await api.post('/v1/agents', {
            name: roleConfig.label,
            type: roleConfig.type,
            description: roleConfig.description,
            config: {
              enable_internet: true,
            },
            profile: roleConfig.profile
          });
          
          if (res.data && res.data.id) {
              setAgentId(res.data.id);
              localStorage.setItem('agent_id', res.data.id);
          }
        } catch (err) {
          console.error('Failed to create agent:', err);
        }
      } else {
        // Update existing agent
        try {
             await api.put(`/v1/agents/${agentId}`, {
                name: roleConfig.label,
                type: roleConfig.type,
                description: roleConfig.description,
                config: {
                    enable_internet: true,
                },
                profile: roleConfig.profile
            });
        } catch (err) {
            console.error('Failed to update agent:', err);
        }
      }
    };

    if (open) {
        initOrUpdateAgent();
    }
  }, [open, agentId, selectedRole]);

  const handleRoleChange = (event: SelectChangeEvent) => {
    setSelectedRole(event.target.value);
  };

  const handleSend = async () => {
    if (!inputValue.trim() || !agentId) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      role: 'user',
      content: inputValue,
      timestamp: Date.now()
    };

    setMessages(prev => [...prev, userMessage]);
    setInputValue('');
    setIsLoading(true);

    const botMessageId = (Date.now() + 1).toString();
    setMessages(prev => [...prev, {
        id: botMessageId,
        role: 'assistant',
        content: '',
        timestamp: Date.now()
    }]);

    abortControllerRef.current = new AbortController();

    try {
        const token = localStorage.getItem('token');
        const response = await fetch('/api/v1/agents/chat', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': token || ''
            },
            body: JSON.stringify({
                agent_id: agentId,
                message: userMessage.content
            }),
            signal: abortControllerRef.current.signal
        });

        if (!response.ok) {
            throw new Error(response.statusText);
        }

        if (!response.body) return;

        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        let botContent = '';
        bufferRef.current = '';

        while (true) {
            const { done, value } = await reader.read();
            if (done) break;

            const chunk = decoder.decode(value, { stream: true });
            bufferRef.current += chunk;
            
            const lines = bufferRef.current.split('\n');
            // Keep the last line in buffer if it's potentially incomplete
            // However, we can't easily know if the last line is complete unless it ends with \n.
            // split('\n') returns an array. If the string ended with \n, the last element is empty string.
            // If it didn't end with \n, the last element is the incomplete line.
            
            // Example: "data: hello\n" -> ["data: hello", ""]
            // Example: "data: hel" -> ["data: hel"]
            
            bufferRef.current = lines.pop() || '';

            for (const line of lines) {
                if (line.trim() === '') continue;

                if (line.startsWith('event: error')) {
                    // Handle error event if needed
                    console.error('Stream error event');
                } else if (line.startsWith('data:')) {
                    // Spec allows optional space after colon
                    let data = line.slice(5);
                    if (data.startsWith(' ')) {
                        data = data.slice(1);
                    }
                    
                    if (data) {
                        botContent += data;
                        setMessages(prev => prev.map(msg => 
                            msg.id === botMessageId 
                                ? { ...msg, content: botContent }
                                : msg
                        ));
                    }
                }
            }
        }

    } catch (error: any) {
        if (error.name === 'AbortError') return;
        console.error('Failed to send message', error);
        setMessages(prev => prev.map(msg => 
            msg.id === botMessageId 
                ? { ...msg, content: 'Sorry, I encountered an error. Please try again.' }
                : msg
        ));
    } finally {
        setIsLoading(false);
        abortControllerRef.current = null;
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <Drawer
      anchor="right"
      open={open}
      onClose={onClose}
      PaperProps={{
        sx: { width: { xs: '100%', sm: 400 }, display: 'flex', flexDirection: 'column' }
      }}
      sx={{ zIndex: (theme) => theme.zIndex.drawer + 2 }} // Ensure it's above other drawers
    >
        {/* Header */}
        <Box sx={{ 
          p: 2, 
          borderBottom: 1, 
          borderColor: 'divider', 
          display: 'flex', 
          alignItems: 'center', 
          justifyContent: 'space-between',
          bgcolor: 'primary.main',
          color: 'primary.contrastText'
        }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <BotIcon />
            <Typography variant="h6">AI Assistant</Typography>
          </Box>
          <IconButton onClick={onClose} color="inherit">
            <CloseIcon />
          </IconButton>
        </Box>

        {/* Role Selector */}
        <Box sx={{ px: 2, py: 1, bgcolor: 'background.paper', borderBottom: 1, borderColor: 'divider' }}>
            <FormControl fullWidth size="small">
                {/* <InputLabel id="role-select-label">Role</InputLabel> */}
                <Select
                    labelId="role-select-label"
                    value={selectedRole}
                    onChange={handleRoleChange}
                    displayEmpty
                    inputProps={{ 'aria-label': 'Select Role' }}
                >
                    {ROLES.map((role) => (
                        <MenuItem key={role.id} value={role.id}>
                            {role.label}
                        </MenuItem>
                    ))}
                </Select>
            </FormControl>
        </Box>

        {/* Messages Area */}
        <Box sx={{ flexGrow: 1, overflow: 'auto', p: 2, bgcolor: 'background.default' }}>
          <List>
            {messages.map((msg) => (
              <ListItem 
                key={msg.id} 
                sx={{ 
                  flexDirection: 'column', 
                  alignItems: msg.role === 'user' ? 'flex-end' : 'flex-start',
                  pb: 2
                }}
              >
                <Box sx={{ 
                  display: 'flex', 
                  gap: 1, 
                  flexDirection: msg.role === 'user' ? 'row-reverse' : 'row',
                  maxWidth: '85%'
                }}>
                  <Paper 
                    elevation={1} 
                    sx={{ 
                      p: 1.5, 
                      bgcolor: msg.role === 'user' ? 'primary.light' : 'background.paper',
                      color: msg.role === 'user' ? 'primary.contrastText' : 'text.primary',
                      borderRadius: 2
                    }}
                  >
                    <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>
                      {msg.content}
                    </Typography>
                  </Paper>
                </Box>
                <Typography variant="caption" color="text.secondary" sx={{ mt: 0.5, px: 1 }}>
                    {msg.role === 'assistant' ? 'AI Agent' : 'You'}
                </Typography>
              </ListItem>
            ))}
            {isLoading && !messages.find(m => m.content === '' && m.role === 'assistant') && (
                 <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2 }}>
                    <CircularProgress size={20} />
                 </Box>
            )}
            <div ref={messagesEndRef} />
          </List>
        </Box>

        {/* Input Area */}
        <Box sx={{ p: 2, borderTop: 1, borderColor: 'divider', bgcolor: 'background.paper' }}>
          <TextField
            fullWidth
            multiline
            maxRows={4}
            placeholder={agentId ? "Type a message..." : "Connecting to agent..."}
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyPress={handleKeyPress}
            disabled={isLoading || !agentId}
            InputProps={{
              endAdornment: (
                <InputAdornment position="end">
                  <IconButton 
                    onClick={handleSend} 
                    disabled={!inputValue.trim() || isLoading || !agentId}
                    color="primary"
                  >
                    <SendIcon />
                  </IconButton>
                </InputAdornment>
              ),
            }}
          />
        </Box>
    </Drawer>
  );
};

export default AgentChatSidebar;
