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
  CircularProgress
} from '@mui/material';
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
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const abortControllerRef = useRef<AbortController | null>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  // Ensure an agent exists
  useEffect(() => {
    const initAgent = async () => {
      if (agentId) return;
      
      try {
        const res = await api.post('/v1/agents', {
          name: 'OmniRead Assistant',
          type: 'reader',
          description: 'Default assistant for reading help',
          config: {
            enable_internet: true,
          },
          profile: {
            avatar: '',
            language: 'Chinese',
            bio: 'A helpful reading assistant.'
          }
        });
        
        // Backend returns { code: 0, msg: "success", data: { id: "..." } }
        if (res.data && res.data.id) {
            setAgentId(res.data.id);
            localStorage.setItem('agent_id', res.data.id);
        }
      } catch (err) {
        console.error('Failed to create agent:', err);
      }
    };

    if (open) {
        initAgent();
    }
  }, [open, agentId]);

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

        while (true) {
            const { done, value } = await reader.read();
            if (done) break;

            const chunk = decoder.decode(value, { stream: true });
            const lines = chunk.split('\n');

            for (const line of lines) {
                if (line.startsWith('event: error')) {
                    // Handle error event if needed
                    console.error('Stream error event');
                } else if (line.startsWith('data: ')) {
                    const data = line.slice(6);
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
