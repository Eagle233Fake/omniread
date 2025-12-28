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
  CircularProgress
} from '@mui/material';
import { 
  ArrowBack as ArrowBackIcon, 
  Settings as SettingsIcon,
  NavigateBefore as PrevIcon,
  NavigateNext as NextIcon
} from '@mui/icons-material';
import { ReactReader } from 'react-reader';
import { Document, Page, pdfjs } from 'react-pdf';
import { useQuery, useMutation } from '@tanstack/react-query';
import { supabase } from '../../api/supabase';
import { useAuth } from '../../context/AuthContext';
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
  const { user } = useAuth();
  
  const [location, setLocation] = useState<string | number>(0);
  const [fontSize, setFontSize] = useState(100);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [numPages, setNumPages] = useState<number>(0);
  const [pageNumber, setPageNumber] = useState<number>(1);
  const [fontFamily, setFontFamily] = useState('Roboto');
  const [signedUrl, setSignedUrl] = useState<string | null>(null);
  const startTimeRef = useRef(0);
  const [rendition, setRendition] = useState<any>(null);

  // Fetch book details
  const { data: book, isLoading } = useQuery<Book>({
    queryKey: ['book', id],
    queryFn: async () => {
      const { data, error } = await supabase
        .from('books')
        .select('*')
        .eq('id', id)
        .single();
      if (error) throw error;
      return data;
    },
    enabled: !!id
  });

  // Get Signed URL
  useEffect(() => {
    if (book?.file_url) {
      supabase.storage.from('book-files').createSignedUrl(book.file_url, 3600)
        .then(({ data }) => setSignedUrl(data?.signedUrl || null));
    }
  }, [book]);

  // Fetch user preferences (from Profile)
  useEffect(() => {
    if (user && id) {
        supabase.from('profiles').select('preferences').eq('id', user.id).single()
        .then(({ data }) => {
             const prefs = data?.preferences as any;
             if (prefs) {
                 if (prefs.font_size) setFontSize(prefs.font_size);
                 if (prefs.font_family) setFontFamily(prefs.font_family);
             }
        });
    }
  }, [user, id]);

  // Fetch progress
  const { data: progress } = useQuery({
    queryKey: ['progress', id],
    queryFn: async () => {
      const { data } = await supabase
        .from('reading_progress')
        .select('*')
        .eq('book_id', id)
        .eq('user_id', user?.id)
        .maybeSingle(); // maybeSingle returns null if not found
      return data;
    },
    enabled: !!id && !!user
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
      if (!user || !id) return;
      return supabase.from('reading_progress').upsert({
        user_id: user.id,
        book_id: id,
        current_loc: data.current_loc,
        progress: data.progress,
        status: data.progress >= 99 ? 'finished' : 'reading',
        updated_at: new Date().toISOString(),
      }, { onConflict: 'user_id,book_id' });
    }
  });

  // Save Session on Unmount
  useEffect(() => {
    // Reset timer on mount
    startTimeRef.current = Date.now();
    
    return () => {
        const endTime = Date.now();
        const duration = Math.round((endTime - startTimeRef.current) / 1000);
        // Only save if read for more than 5 seconds
        if (duration > 5 && user && id) {
             supabase.from('reading_sessions').insert({
                 user_id: user.id,
                 book_id: id,
                 start_time: new Date(startTimeRef.current).toISOString(),
                 end_time: new Date(endTime).toISOString(),
                 duration: duration
             });
        }
    };
  }, [user, id]); // Re-run if user/id changes, effectively saving previous session

  const debouncedSave = useDebounce((loc: string | number, prog: number) => {
    saveProgressMutation.mutate({
      current_loc: loc.toString(),
      progress: prog
    });
  }, 1000);

  // Update EPUB font size
  useEffect(() => {
    if (rendition) {
      rendition.themes.fontSize(`${fontSize}%`);
      rendition.themes.font(fontFamily);
    }
  }, [fontSize, fontFamily, rendition]);

  const onLocationChange = (loc: string | number) => {
    setLocation(loc);
    debouncedSave(loc, 0); 
  };

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
      <Toolbar />

      <Box sx={{ flexGrow: 1, position: 'relative', overflow: 'hidden' }}>
        {book?.format === 'epub' && signedUrl && (
          <Box sx={{ height: '100%' }}>
            <ReactReader
              url={signedUrl}
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

        {book?.format === 'pdf' && signedUrl && (
          <Box sx={{ 
            height: '100%', 
            display: 'flex', 
            justifyContent: 'center', 
            overflow: 'auto',
            bgcolor: 'grey.100',
            pt: 2
          }}>
            <Document
              file={signedUrl}
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
