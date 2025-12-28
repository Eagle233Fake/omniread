import React, { useState } from 'react';
import { Box, Button, TextField, Typography, Container, Paper, Link, Alert } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { supabase } from '../../api/supabase';

const Register: React.FC = () => {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const { data, error } = await supabase.auth.signUp({
        email,
        password,
        options: {
          data: {
            username,
            avatar_url: `https://api.dicebear.com/7.x/avataaars/svg?seed=${username}`, // Generate default avatar
          },
        },
      });

      if (error) {
        setError(error.message);
      } else if (data.user) {
        // Auto sign-in or redirect (Supabase default is confirm email unless disabled)
        // Assuming email confirmation is disabled for dev, or handling "check email"
        if (data.session) {
            navigate('/');
        } else {
            setError('Please check your email to verify your account.');
        }
      }
    } catch (err) {
      setError('Registration failed.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container component="main" maxWidth="xs">
      <Box
        sx={{
          marginTop: 8,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <Typography component="h1" variant="h4" sx={{ mb: 4, color: 'primary.main' }}>
          OmniRead
        </Typography>
        <Paper sx={{ p: 4, width: '100%', display: 'flex', flexDirection: 'column', gap: 2 }}>
          <Typography component="h2" variant="h5" sx={{ mb: 2 }}>
            Sign Up
          </Typography>
          <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              required
              fullWidth
              label="Username"
              autoComplete="username"
              autoFocus
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
             <TextField
              required
              fullWidth
              label="Email Address"
              autoComplete="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
            <TextField
              required
              fullWidth
              label="Password"
              type="password"
              autoComplete="new-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            {error && (
              <Alert severity="info">{error}</Alert>
            )}
            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 2, mb: 2 }}
              disabled={loading}
            >
              {loading ? 'Signing up...' : 'Sign Up'}
            </Button>
            <Link href="/login" variant="body2" sx={{ textAlign: 'center' }}>
              {"Already have an account? Sign In"}
            </Link>
          </Box>
        </Paper>
      </Box>
    </Container>
  );
};

export default Register;
