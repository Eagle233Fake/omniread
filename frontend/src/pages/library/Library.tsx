import React, { useState } from 'react';
import { 
  Box, 
  Grid, 
  Card, 
  CardContent, 
  CardMedia, 
  Typography, 
  Fab, 
  Dialog, 
  DialogTitle, 
  DialogContent, 
  DialogActions, 
  Button, 
  TextField,
  CardActionArea,
  Chip,
  IconButton,
  Menu,
  MenuItem,
  ListItemIcon
} from '@mui/material';
import { 
  Add as AddIcon, 
  CloudUpload as CloudUploadIcon,
  MoreVert as MoreVertIcon,
  Edit as EditIcon
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '../../api/client';

interface Book {
  id: string;
  title: string;
  author: string;
  cover_url: string;
  format: string;
  description: string;
  created_at: string;
}

import { useNavigate } from 'react-router-dom';

const Library: React.FC = () => {
  const [open, setOpen] = useState(false);
  const [file, setFile] = useState<File | null>(null);
  const [title, setTitle] = useState('');
  const [author, setAuthor] = useState('');
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const [editOpen, setEditOpen] = useState(false);
  const [editingBook, setEditingBook] = useState<Book | null>(null);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedBookId, setSelectedBookId] = useState<string | null>(null);

  const { data: books, isLoading } = useQuery<Book[]>({
    queryKey: ['books'],
    queryFn: async () => {
      const res: any = await api.get('/books');
      return res.data || [];
    }
  });

  const uploadMutation = useMutation({
    mutationFn: async (formData: FormData) => {
      return api.post('/books/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['books'] });
      handleClose();
    }
  });

  const updateMutation = useMutation({
    mutationFn: async (data: { id: string, title: string, author: string, description: string }) => {
      return api.put(`/books/${data.id}`, {
        title: data.title,
        author: data.author,
        description: data.description
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['books'] });
      handleEditClose();
    }
  });

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0]);
    }
  };

  const handleEditClick = (event: React.MouseEvent<HTMLElement>, book: Book) => {
    event.stopPropagation();
    setAnchorEl(event.currentTarget);
    setSelectedBookId(book.id);
    setEditingBook(book);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
    setSelectedBookId(null);
  };

  const handleEditOpen = () => {
    handleMenuClose();
    if (editingBook) {
      setTitle(editingBook.title);
      setAuthor(editingBook.author);
      setEditOpen(true);
    }
  };

  const handleEditClose = () => {
    setEditOpen(false);
    setEditingBook(null);
    setTitle('');
    setAuthor('');
  };

  const handleUpdate = () => {
    if (editingBook) {
      updateMutation.mutate({
        id: editingBook.id,
        title,
        author,
        description: editingBook.description || '' // Should load description too if possible
      });
    }
  };

  const handleUpload = () => {
    if (!file) return;
    const formData = new FormData();
    formData.append('file', file);
    if (title) formData.append('title', title);
    if (author) formData.append('author', author);
    uploadMutation.mutate(formData);
  };

  const handleClose = () => {
    setOpen(false);
    setFile(null);
    setTitle('');
    setAuthor('');
  };

  return (
    <Box>
      <Typography variant="h4" sx={{ mb: 4, color: 'text.primary' }}>
        My Library
      </Typography>

      {isLoading ? (
        <Typography>Loading...</Typography>
      ) : (
        <Grid container spacing={3}>
          {books?.length === 0 && (
            <Grid item xs={12}>
              <Typography variant="body1" align="center" sx={{ width: '100%', mt: 4, color: 'text.secondary' }}>
                No books found in your library.
              </Typography>
            </Grid>
          )}
          {books?.map((book) => (
            <Grid item xs={12} sm={6} md={4} lg={3} key={book.id}>
              <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
                <CardActionArea 
                  onClick={() => navigate(`/read/${book.id}`)}
                  sx={{ flexGrow: 1, display: 'flex', flexDirection: 'column', alignItems: 'stretch' }}
                >
                  <CardMedia
                    component="img"
                    height="200"
                    image={book.cover_url || '/api/uploads/default_cover.svg'}
                    alt={book.title}
                    sx={{ objectFit: 'cover', aspectRatio: '3/4' }}
                    onError={(e: any) => {
                      e.target.onerror = null; 
                      e.target.src = '/api/uploads/default_cover.svg';
                    }}
                  />
                  <CardContent sx={{ position: 'relative' }}>
                    <Box sx={{ position: 'absolute', top: 8, right: 8 }}>
                      <IconButton 
                        size="small"
                        onClick={(e) => handleEditClick(e, book)}
                      >
                        <MoreVertIcon />
                      </IconButton>
                    </Box>
                    <Typography gutterBottom variant="h6" component="div" noWrap sx={{ pr: 3 }}>
                      {book.title}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {book.author || 'Unknown Author'}
                    </Typography>
                     <Box sx={{ mt: 1 }}>
                      <Chip label={book.format} size="small" color="primary" variant="outlined" />
                    </Box>
                  </CardContent>
                </CardActionArea>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}

      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
        onClick={(e) => e.stopPropagation()}
      >
        <MenuItem onClick={handleEditOpen}>
          <ListItemIcon>
            <EditIcon fontSize="small" />
          </ListItemIcon>
          Edit
        </MenuItem>
      </Menu>

      <Dialog open={editOpen} onClose={handleEditClose} maxWidth="sm" fullWidth>
        <DialogTitle>Edit Book Details</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
            <TextField
              label="Title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              fullWidth
            />
            <TextField
              label="Author"
              value={author}
              onChange={(e) => setAuthor(e.target.value)}
              fullWidth
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleEditClose}>Cancel</Button>
          <Button 
            onClick={handleUpdate} 
            variant="contained"
            disabled={updateMutation.isPending}
          >
            {updateMutation.isPending ? 'Saving...' : 'Save'}
          </Button>
        </DialogActions>
      </Dialog>

      <Fab 
        color="primary" 
        aria-label="add" 
        sx={{ position: 'fixed', bottom: 32, right: 32 }}
        onClick={() => setOpen(true)}
      >
        <AddIcon />
      </Fab>

      <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
        <DialogTitle>Add New Book</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
            <Button
              variant="outlined"
              component="label"
              startIcon={<CloudUploadIcon />}
              fullWidth
            >
              {file ? file.name : 'Select PDF/EPUB File'}
              <input
                type="file"
                hidden
                accept=".pdf,.epub"
                onChange={handleFileChange}
              />
            </Button>
            <TextField
              label="Title (Optional)"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              fullWidth
            />
            <TextField
              label="Author (Optional)"
              value={author}
              onChange={(e) => setAuthor(e.target.value)}
              fullWidth
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Cancel</Button>
          <Button 
            onClick={handleUpload} 
            variant="contained" 
            disabled={!file || uploadMutation.isPending}
          >
            {uploadMutation.isPending ? 'Uploading...' : 'Upload'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default Library;
