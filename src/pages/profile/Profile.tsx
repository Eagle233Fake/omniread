import React, { useState, useEffect } from 'react';
import { 
  Box, Typography, Paper, TextField, Button, Avatar, Grid, Tabs, Tab,
  Alert, Snackbar, Select, MenuItem, FormControl, InputLabel
} from '@mui/material';
import { useAuth } from '../../context/AuthContext';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { supabase } from '../../api/supabase';

const Profile: React.FC = () => {
  const { user } = useAuth();
  const [tab, setTab] = useState(0);
  const [msg, setMsg] = useState({ type: '', text: '' });
  const queryClient = useQueryClient();
  
  // Profile Form State
  const [profile, setProfile] = useState({
    username: '',
    email: '',
    phone: '',
    bio: ''
  });

  // Password Form State
  const [password, setPassword] = useState({
    old: '',
    new: '',
    confirm: ''
  });

  // Preferences State
  const [preferences, setPreferences] = useState({
    font_family: 'Roboto',
    font_size: 100
  });

  // Fetch current user details
  const { data: userData } = useQuery({
    queryKey: ['profile', user?.id],
    queryFn: async () => {
      if (!user) return null;
      const { data, error } = await supabase
        .from('profiles')
        .select('*')
        .eq('id', user.id)
        .single();
      if (error) throw error;
      return data;
    },
    enabled: !!user
  });

  useEffect(() => {
    if (userData) {
      setProfile({
        username: userData.username || '',
        email: user?.email || '',
        phone: userData.phone || '',
        bio: userData.bio || ''
      });
      if (userData.preferences) {
        setPreferences(userData.preferences as any);
      }
    } else if (user) {
      setProfile(prev => ({
        ...prev,
        username: user.username,
        email: user.email
      }));
    }
  }, [userData, user]);

  const updateProfileMutation = useMutation({
    mutationFn: async (data: any) => {
      if (!user) return;
      const { error } = await supabase.from('profiles').update({
        username: data.username,
        phone: data.phone,
        bio: data.bio
      }).eq('id', user.id);
      if (error) throw error;
    },
    onSuccess: () => {
      setMsg({ type: 'success', text: 'Profile updated successfully' });
      queryClient.invalidateQueries({ queryKey: ['profile'] });
    },
    onError: () => setMsg({ type: 'error', text: 'Failed to update profile' })
  });

  const changePasswordMutation = useMutation({
    mutationFn: async (data: any) => {
      const { error } = await supabase.auth.updateUser({
        password: data.new_password
      });
      if (error) throw error;
    },
    onSuccess: () => {
      setMsg({ type: 'success', text: 'Password changed successfully' });
      setPassword({ old: '', new: '', confirm: '' });
    },
    onError: (err: any) => setMsg({ type: 'error', text: err.message || 'Failed to change password' })
  });

  const updatePreferencesMutation = useMutation({
    mutationFn: async (data: any) => {
      if (!user) return;
      const { error } = await supabase.from('profiles').update({
        preferences: data
      }).eq('id', user.id);
      if (error) throw error;
    },
    onSuccess: () => {
      setMsg({ type: 'success', text: 'Preferences saved' });
      queryClient.invalidateQueries({ queryKey: ['profile'] });
    },
    onError: () => setMsg({ type: 'error', text: 'Failed to save preferences' })
  });

  const handleProfileSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    updateProfileMutation.mutate(profile);
  };

  const handlePasswordSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (password.new !== password.confirm) {
      setMsg({ type: 'error', text: 'Passwords do not match' });
      return;
    }
    changePasswordMutation.mutate({
      new_password: password.new
    });
  };

  const handlePreferencesSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    updatePreferencesMutation.mutate(preferences);
  };

  return (
    <Box maxWidth="md" mx="auto">
      <Typography variant="h4" sx={{ mb: 4 }}>Account Settings</Typography>
      
      <Paper sx={{ mb: 4 }}>
        <Tabs value={tab} onChange={(_, v) => setTab(v)} sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tab label="Profile" />
          <Tab label="Security" />
          <Tab label="Preferences" />
        </Tabs>

        <Box p={3}>
          {tab === 0 && (
            <form onSubmit={handleProfileSubmit}>
              <Grid container spacing={2} alignItems="center" sx={{ mb: 4 }}>
                <Grid size={{ xs: 12 }} display="flex" justifyContent="center">
                  <Box position="relative">
                    <Avatar 
                      sx={{ width: 100, height: 100, fontSize: 40 }}
                      src={user?.avatar}
                    >
                      {profile.username?.[0]?.toUpperCase()}
                    </Avatar>
                  </Box>
                </Grid>
                <Grid size={{ xs: 12 }}>
                  <TextField 
                    fullWidth label="Nickname / Username" 
                    value={profile.username}
                    onChange={e => setProfile({...profile, username: e.target.value})}
                  />
                </Grid>
                <Grid size={{ xs: 12, sm: 6 }}>
                  <TextField
                    fullWidth label="Email" 
                    value={profile.email}
                    disabled
                  />
                </Grid>
                <Grid size={{ xs: 12, sm: 6 }}>
                  <TextField 
                    fullWidth label="Phone" 
                    value={profile.phone}
                    onChange={e => setProfile({...profile, phone: e.target.value})}
                  />
                </Grid>
                <Grid size={{ xs: 12 }}>
                  <TextField 
                    fullWidth multiline rows={3} label="Bio" 
                    value={profile.bio}
                    onChange={e => setProfile({...profile, bio: e.target.value})}
                  />
                </Grid>
                <Grid size={{ xs: 12 }}>
                  <Button variant="contained" type="submit" disabled={updateProfileMutation.isPending}>
                    Save Changes
                  </Button>
                </Grid>
              </Grid>
            </form>
          )}

          {tab === 1 && (
            <form onSubmit={handlePasswordSubmit}>
              <Grid container spacing={3} maxWidth="sm">
                <Grid size={{ xs: 12 }}>
                   <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                     Enter your new password below.
                   </Typography>
                </Grid>
                <Grid size={{ xs: 12 }}>
                  <TextField 
                    fullWidth type="password" label="New Password" 
                    value={password.new}
                    onChange={e => setPassword({...password, new: e.target.value})}
                  />
                </Grid>
                <Grid size={{ xs: 12 }}>
                  <TextField 
                    fullWidth type="password" label="Confirm New Password" 
                    value={password.confirm}
                    onChange={e => setPassword({...password, confirm: e.target.value})}
                  />
                </Grid>
                <Grid size={{ xs: 12 }}>
                  <Button variant="contained" type="submit" color="error" disabled={changePasswordMutation.isPending}>
                    Change Password
                  </Button>
                </Grid>
              </Grid>
            </form>
          )}

          {tab === 2 && (
            <form onSubmit={handlePreferencesSubmit}>
              <Grid container spacing={3} maxWidth="sm">
                <Grid size={{ xs: 12, sm: 6 }}>
                  <FormControl fullWidth>
                    <InputLabel>Reading Font</InputLabel>
                    <Select
                      value={preferences.font_family}
                      label="Reading Font"
                      onChange={e => setPreferences({...preferences, font_family: e.target.value})}
                    >
                      <MenuItem value="Roboto">Roboto (Default)</MenuItem>
                      <MenuItem value="Merriweather">Merriweather (Serif)</MenuItem>
                      <MenuItem value="Open Sans">Open Sans</MenuItem>
                      <MenuItem value="Lora">Lora</MenuItem>
                      <MenuItem value="Source Code Pro">Monospace</MenuItem>
                    </Select>
                  </FormControl>
                </Grid>
                <Grid size={{ xs: 12 }}>
                  <TextField 
                    fullWidth type="number" label="Base Font Size (%)" 
                    value={preferences.font_size}
                    onChange={e => setPreferences({...preferences, font_size: parseInt(e.target.value)})}
                    helperText="Default size percentage for reader (e.g. 100)"
                  />
                </Grid>
                <Grid size={{ xs: 12 }}>
                  <Button variant="contained" type="submit" disabled={updatePreferencesMutation.isPending}>
                    Save Preferences
                  </Button>
                </Grid>
              </Grid>
            </form>
          )}
        </Box>
      </Paper>

      <Snackbar 
        open={!!msg.text} 
        autoHideDuration={6000} 
        onClose={() => setMsg({ type: '', text: '' })}
      >
        <Alert severity={msg.type as any} onClose={() => setMsg({ type: '', text: '' })}>
          {msg.text}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default Profile;
