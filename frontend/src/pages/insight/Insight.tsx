import React from 'react';
import { Box, Grid, Card, CardContent, Typography, useTheme } from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Cell } from 'recharts';
import api from '../../api/client';
import { format, parseISO } from 'date-fns';

interface DailyStat {
  date: string;
  duration: number;
}

interface SummaryData {
  total_reading_time: number;
  total_books_read: number;
  current_streak: number;
  daily_stats: DailyStat[];
}

const Insight: React.FC = () => {
  const theme = useTheme();
  
  const { data: summary, isLoading } = useQuery<SummaryData>({
    queryKey: ['insight'],
    queryFn: async () => {
      const res: any = await api.get('/insight/summary');
      return res.data;
    }
  });

  if (isLoading) return <Typography>Loading...</Typography>;
  // if (!summary) return <Typography>No data available</Typography>;
  // Fallback to empty state if data is missing or error
  const safeSummary = summary || {
    total_reading_time: 0,
    total_books_read: 0,
    current_streak: 0,
    daily_stats: []
  };

  // Format daily stats for chart
  const chartData = safeSummary.daily_stats?.map(item => ({
    date: format(parseISO(item.date), 'MMM dd'),
    duration: Math.round(item.duration / 60) // Convert seconds to minutes
  })) || [];

  const StatCard = ({ title, value, unit }: { title: string, value: number | string, unit?: string }) => (
    <Card sx={{ height: '100%', bgcolor: 'primary.light', color: 'primary.dark' }}>
      <CardContent>
        <Typography variant="subtitle2" sx={{ opacity: 0.8 }}>
          {title}
        </Typography>
        <Typography variant="h3" sx={{ fontWeight: 'bold', my: 1 }}>
          {value}
        </Typography>
        {unit && <Typography variant="body2">{unit}</Typography>}
      </CardContent>
    </Card>
  );

  return (
    <Box>
      <Typography variant="h4" sx={{ mb: 4 }}>Reading Insights</Typography>
      
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} md={4}>
          <StatCard 
            title="Total Reading Time" 
            value={Math.round(safeSummary.total_reading_time / 3600)} 
            unit="Hours" 
          />
        </Grid>
        <Grid item xs={12} md={4}>
          <StatCard 
            title="Books Completed" 
            value={safeSummary.total_books_read} 
            unit="Books" 
          />
        </Grid>
        <Grid item xs={12} md={4}>
          <StatCard 
            title="Current Streak" 
            value={safeSummary.current_streak} 
            unit="Days" 
          />
        </Grid>
      </Grid>

      <Card sx={{ p: 3 }}>
        <Typography variant="h6" sx={{ mb: 3 }}>Daily Reading Activity (Last 30 Days)</Typography>
        <Box sx={{ height: 300, width: '100%' }}>
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={chartData}>
              <XAxis 
                dataKey="date" 
                tick={{ fontSize: 12 }} 
                tickLine={false}
                axisLine={false}
              />
              <YAxis 
                hide 
              />
              <Tooltip 
                cursor={{ fill: 'transparent' }}
                contentStyle={{ 
                  backgroundColor: theme.palette.background.paper,
                  borderRadius: 8,
                  border: 'none',
                  boxShadow: theme.shadows[2]
                }}
              />
              <Bar dataKey="duration" radius={[4, 4, 0, 0]}>
                {chartData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={theme.palette.primary.main} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </Box>
      </Card>
    </Box>
  );
};

export default Insight;
