import React, { useState } from 'react';
import { 
  Box, 
  Drawer, 
  List, 
  ListItem, 
  ListItemButton, 
  ListItemIcon, 
  ListItemText, 
  AppBar, 
  Toolbar, 
  Typography, 
  IconButton, 
  Avatar, 
  Menu, 
  MenuItem,
  useTheme,
  useMediaQuery
} from '@mui/material';
import { 
  LibraryBooks as LibraryIcon, 
  Analytics as InsightIcon, 
  Menu as MenuIcon,
  Logout as LogoutIcon
} from '@mui/icons-material';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

import { 
  MenuOpen as MenuOpenIcon,
  ChevronLeft as ChevronLeftIcon
} from '@mui/icons-material';

const drawerWidth = 240;
const railWidth = 80;

const Layout: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [mobileOpen, setMobileOpen] = useState(false);
  const [desktopOpen, setDesktopOpen] = useState(false); // New state for desktop drawer

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    handleClose();
    logout();
    navigate('/login');
  };

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const handleDesktopToggle = () => {
    setDesktopOpen(!desktopOpen);
  };

  const menuItems = [
    { text: 'Library', icon: <LibraryIcon />, path: '/library' },
    { text: 'Insight', icon: <InsightIcon />, path: '/insight' },
  ];

  // Helper to determine if we are in "open" state (mobile drawer open or desktop drawer expanded)
  const isDrawerOpen = isMobile ? mobileOpen : desktopOpen;
  const currentWidth = isMobile ? drawerWidth : (desktopOpen ? drawerWidth : railWidth);

  const drawer = (
    <Box sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <Toolbar sx={{ justifyContent: isDrawerOpen ? 'space-between' : 'center', py: 2, px: isDrawerOpen ? 2 : 1 }}>
        {isDrawerOpen && (
          <Typography variant="h6" noWrap component="div" sx={{ color: 'primary.main', fontWeight: 'bold' }}>
             OmniRead
          </Typography>
        )}
        {!isMobile && (
          <IconButton onClick={handleDesktopToggle}>
             {desktopOpen ? <ChevronLeftIcon /> : <MenuOpenIcon />}
          </IconButton>
        )}
      </Toolbar>
      <List sx={{ flexGrow: 1, px: 1 }}>
        {menuItems.map((item) => (
          <ListItem key={item.text} disablePadding sx={{ display: 'block', mb: 1 }}>
            <ListItemButton
              selected={location.pathname === item.path}
              onClick={() => {
                navigate(item.path);
                if (isMobile) setMobileOpen(false);
              }}
              sx={{
                minHeight: 48,
                justifyContent: isDrawerOpen ? 'initial' : 'center',
                px: 2.5,
                borderRadius: 4, // M3 shape
                flexDirection: isDrawerOpen ? 'row' : 'column',
                '&.Mui-selected': {
                  bgcolor: 'secondary.light',
                  color: 'secondary.dark',
                  '&:hover': {
                    bgcolor: 'secondary.light',
                  },
                  '& .MuiListItemIcon-root': {
                    color: 'secondary.dark',
                  }
                },
              }}
            >
              <ListItemIcon
                sx={{
                  minWidth: 0,
                  mr: isDrawerOpen ? 3 : 0,
                  justifyContent: 'center',
                  mb: isDrawerOpen ? 0 : 0.5,
                }}
              >
                {item.icon}
              </ListItemIcon>
              <ListItemText 
                primary={item.text} 
                sx={{ 
                  opacity: isDrawerOpen ? 1 : 1, 
                  textAlign: isDrawerOpen ? 'left' : 'center',
                  display: isDrawerOpen ? 'block' : 'none', // Hide text when collapsed if strict rail
                  '& .MuiTypography-root': {
                    fontSize: isDrawerOpen ? '1rem' : '0.75rem',
                    fontWeight: 500
                  }
                }} 
              />
               {/* Show mini text in rail mode if desired, or hide completely */}
               {!isDrawerOpen && (
                 <Typography variant="caption" sx={{ fontSize: '0.65rem', fontWeight: 500, mt: 0.5 }}>
                   {item.text}
                 </Typography>
               )}
            </ListItemButton>
          </ListItem>
        ))}
      </List>
    </Box>
  );

  return (
    <Box sx={{ display: 'flex' }}>
      <AppBar
        position="fixed"
        sx={{
          width: { sm: `calc(100% - ${currentWidth}px)` },
          ml: { sm: `${currentWidth}px` },
          bgcolor: 'background.default',
          color: 'text.primary',
          boxShadow: 'none',
          borderBottom: '1px solid',
          borderColor: 'divider',
          transition: theme.transitions.create(['width', 'margin'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
          }),
        }}
      >
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ mr: 2, display: { sm: 'none' } }}
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1 }}>
            {menuItems.find(i => i.path === location.pathname)?.text || 'OmniRead'}
          </Typography>
          
          <Box>
            <IconButton
              size="large"
              aria-label="account of current user"
              aria-controls="menu-appbar"
              aria-haspopup="true"
              onClick={handleMenu}
              color="inherit"
            >
              <Avatar sx={{ bgcolor: 'primary.main' }}>{user?.username?.charAt(0).toUpperCase()}</Avatar>
            </IconButton>
            <Menu
              id="menu-appbar"
              anchorEl={anchorEl}
              anchorOrigin={{
                vertical: 'top',
                horizontal: 'right',
              }}
              keepMounted
              transformOrigin={{
                vertical: 'top',
                horizontal: 'right',
              }}
              open={Boolean(anchorEl)}
              onClose={handleClose}
            >
              <MenuItem onClick={() => { handleClose(); navigate('/profile'); }}>Profile & Settings</MenuItem>
              <MenuItem onClick={handleLogout}>
                <ListItemIcon>
                  <LogoutIcon fontSize="small" />
                </ListItemIcon>
                Logout
              </MenuItem>
            </Menu>
          </Box>
        </Toolbar>
      </AppBar>

      <Box
        component="nav"
        sx={{ width: { sm: currentWidth }, flexShrink: { sm: 0 }, transition: theme.transitions.create('width', {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
          }) 
        }}
      >
        {/* Mobile Drawer */}
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{
            keepMounted: true, // Better open performance on mobile.
          }}
          sx={{
            display: { xs: 'block', sm: 'none' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
          }}
        >
          {drawer}
        </Drawer>

        {/* Desktop Rail/Drawer */}
        <Drawer
          variant="permanent"
          sx={{
            display: { xs: 'none', sm: 'block' },
            '& .MuiDrawer-paper': { 
                boxSizing: 'border-box', 
                width: currentWidth,
                borderRight: '1px solid',
                borderColor: 'divider',
                bgcolor: 'background.paper',
                overflowX: 'hidden',
                transition: theme.transitions.create('width', {
                  easing: theme.transitions.easing.sharp,
                  duration: theme.transitions.duration.leavingScreen,
                }),
            },
          }}
          open
        >
          {drawer}
        </Drawer>
      </Box>

      <Box
        component="main"
        sx={{ 
          flexGrow: 1, 
          p: 3, 
          width: { sm: `calc(100% - ${currentWidth}px)` },
          mt: 8,
          minHeight: '100vh',
          bgcolor: 'background.default',
          transition: theme.transitions.create(['width', 'margin'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
          }),
        }}
      >
        <Outlet />
      </Box>
    </Box>
  );
};

export default Layout;
