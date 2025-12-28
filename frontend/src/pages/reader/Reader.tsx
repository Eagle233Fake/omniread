import React, { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { 
  Box, 
  AppBar, 
  Toolbar, 
  IconButton, 
  Typography, 
  Slider, 
  Menu, 
  MenuItem,
  CircularProgress,
  Fab
} from '@mui/material';
import { 
  ArrowBack as ArrowBackIcon, 
  Settings as SettingsIcon,
  NavigateBefore as PrevIcon,
  NavigateNext as NextIcon
} from '@mui/icons-material';
import { ReactReader } from 'react-reader';
import { Document, Page, pdfjs } from 'react-pdf';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '../../api/client';
import useDebounce from '../../utils/useDebounce';

// Configure PDF.js worker
pdfjs.GlobalWorkerOptions.workerSrc = new URL(
  'pdfjs-dist/build/pdf.worker.min.mjs',
  import.meta.url,
).toString();

interface Book {
  id: string;
  title: string;
  file_url: string;
  format: string;
}

const Reader: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [location, setLocation] = useState<string | number>(0);
  const [fontSize, setFontSize] = useState(100);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [numPages, setNumPages] = useState<number>(0);
  const [pageNumber, setPageNumber] = useState<number>(1);
  
  const [fontFamily, setFontFamily] = useState('Roboto');

  // Fetch book details
  const { data: book, isLoading } = useQuery<Book>({
    queryKey: ['book', id],
    queryFn: async () => {
      const res: any = await api.get(`/books/${id}`);
      return res.data;
    },
    enabled: !!id
  });

  // Fetch user preferences
  const { data: userProfile } = useQuery({
    queryKey: ['profile'],
    queryFn: async () => {
      const res: any = await api.get('/user/profile');
      return res.data;
    }
  });

  useEffect(() => {
    if (userProfile?.preferences) {
      if (userProfile.preferences.font_size) setFontSize(userProfile.preferences.font_size);
      if (userProfile.preferences.font_family) setFontFamily(userProfile.preferences.font_family);
    }
  }, [userProfile]);

  // Fetch progress
  const { data: progress } = useQuery({
    queryKey: ['progress', id],
    queryFn: async () => {
      const res: any = await api.get(`/reading/progress?book_id=${id}`);
      return res.data;
    },
    enabled: !!id
  });

  // Initialize location from progress
  useEffect(() => {
    if (progress) {
      if (book?.format === 'epub') {
        setLocation(progress.current_loc || 0);
      } else if (book?.format === 'pdf') {
        setPageNumber(parseInt(progress.current_loc) || 1);
      }
    }
  }, [progress, book]);

  // Save progress mutation
  const saveProgressMutation = useMutation({
    mutationFn: async (data: { current_loc: string, progress: number }) => {
      return api.post('/reading/progress', {
        book_id: id,
        status: 'reading',
        ...data
      });
    }
  });

  const debouncedSave = useDebounce((loc: string | number, prog: number) => {
    saveProgressMutation.mutate({
      current_loc: loc.toString(),
      progress: prog
    });
  }, 1000);

  const [rendition, setRendition] = useState<any>(null);

  // Update EPUB font size when slider changes
  useEffect(() => {
    if (rendition) {
      rendition.themes.fontSize(`${fontSize}%`);
      rendition.themes.font(fontFamily);
    }
  }, [fontSize, fontFamily, rendition]);

  // EPUB location change
  const onLocationChange = (loc: string | number) => {
    setLocation(loc);
    debouncedSave(loc, 0); 
  };

  // PDF page change
  const onPdfPageChange = (offset: number) => {
    const newPage = Math.min(Math.max(1, pageNumber + offset), numPages);
    setPageNumber(newPage);
    const percentage = (newPage / numPages) * 100;
    debouncedSave(newPage, percentage);
  };

  const handleSettingsOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleSettingsClose = () => {
    setAnchorEl(null);
  };

  if (isLoading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ height: '100vh', display: 'flex', flexDirection: 'column', bgcolor: 'background.default' }}>
      {/* Top Bar */}
      <AppBar position="fixed" sx={{ zIndex: 1201 }}>
        <Toolbar>
          <IconButton edge="start" color="inherit" onClick={() => navigate('/library')} sx={{ mr: 2 }}>
            <ArrowBackIcon />
          </IconButton>
          <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1 }}>
            {book?.title}
          </Typography>
          <IconButton color="inherit" onClick={handleSettingsOpen}>
            <SettingsIcon />
          </IconButton>
          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleSettingsClose}
          >
            <Box sx={{ width: 200, p: 2 }}>
              <Typography gutterBottom>Font Size</Typography>
              <Slider
                value={fontSize}
                onChange={(_, val) => setFontSize(val as number)}
                min={80}
                max={150}
                step={10}
              />
            </Box>
          </Menu>
        </Toolbar>
      </AppBar>
      <Toolbar /> {/* Spacer */}

      {/* Reader Content */}
      <Box sx={{ flexGrow: 1, position: 'relative', overflow: 'hidden' }}>
        {book?.format === 'epub' && (
          <Box sx={{ height: '100%' }}>
            <ReactReader
              url={book.file_url}
              location={location}
              locationChanged={onLocationChange}
              epubOptions={{
                flow: 'paginated',
                manager: 'default',
              }}
              getRendition={(r) => {
                setRendition(r);
                r.themes.fontSize(`${fontSize}%`);
                r.themes.font(fontFamily);
              }}
            />
          </Box>
        )}

        {book?.format === 'pdf' && (
          <Box sx={{ 
            height: '100%', 
            display: 'flex', 
            justifyContent: 'center', 
            overflow: 'auto',
            bgcolor: 'grey.100',
            pt: 2
          }}>
            <Document
              file={book.file_url}
              onLoadSuccess={({ numPages }) => setNumPages(numPages)}
              loading={<CircularProgress />}
            >
              <Page 
                pageNumber={pageNumber} 
                scale={fontSize / 100} 
                renderTextLayer={false} 
                renderAnnotationLayer={false}
                className="pdf-page-shadow"
              />
            </Document>

            {/* PDF Navigation Controls */}
            <Box sx={{ position: 'fixed', bottom: 32, left: '50%', transform: 'translateX(-50%)', display: 'flex', gap: 2, bgcolor: 'background.paper', p: 1, borderRadius: 4, boxShadow: 3 }}>
              <IconButton onClick={() => onPdfPageChange(-1)} disabled={pageNumber <= 1}>
                <PrevIcon />
              </IconButton>
              <Typography variant="body1" sx={{ alignSelf: 'center' }}>
                {pageNumber} / {numPages}
              </Typography>
              <IconButton onClick={() => onPdfPageChange(1)} disabled={pageNumber >= numPages}>
                <NextIcon />
              </IconButton>
            </Box>
          </Box>
        )}
      </Box>
    </Box>
  );
};

export default Reader;
