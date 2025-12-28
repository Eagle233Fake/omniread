import React from 'react';
import { Box, Grid, Card, CardContent, Typography, useTheme } from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Cell } from 'recharts';
import { supabase } from '../../api/supabase';
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

const Insight: React.FC = () => {
  const theme = useTheme();
  
  const { data: summary, isLoading } = useQuery<SummaryData>({
    queryKey: ['insight'],
    queryFn: async () => {
      const { data, error } = await supabase.rpc('get_reading_stats', { days: 30 });
      if (error) throw error;
      
      const rpcData = data as any;
      
      // Calculate Streak
      let streak = 0;
      const stats = rpcData.daily_stats || [];
      // Loop from end (Today) backwards
      // Note: stats are ordered by date ASC from SQL
      let hasBroken = false;
      // Check today first. If today has data, good. If not, check yesterday.
      // If neither, streak is 0.
      
      // We need to handle timezone carefully, but roughly:
      // If stats[last].duration > 0 -> streak++
      // If stats[last] == 0, check stats[last-1]. If > 0, streak starts there.
      
      for (let i = stats.length - 1; i >= 0; i--) {
          if (stats[i].duration > 0) {
              streak++;
              hasBroken = false;
          } else {
              // Allow today to be 0 if we are just checking streak
              if (i === stats.length - 1) {
                  continue;
              }
              break;
          }
      }

      return {
        total_reading_time: rpcData.total_duration,
        total_books_read: rpcData.books_finished,
        current_streak: streak,
        daily_stats: stats
      };
    }
  });

  if (isLoading) return <Typography>Loading...</Typography>;
  
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

  return (
    <Box>
      <Typography variant="h4" sx={{ mb: 4 }}>Reading Insights</Typography>
      
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid size={{ xs: 12, md: 4 }}>
          <StatCard 
            title="Total Reading Time" 
            value={Math.round(safeSummary.total_reading_time / 3600)} 
            unit="Hours" 
          />
        </Grid>
        <Grid size={{ xs: 12, md: 4 }}>
          <StatCard 
            title="Books Completed" 
            value={safeSummary.total_books_read} 
            unit="Books" 
          />
        </Grid>
        <Grid size={{ xs: 12, md: 4 }}>
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
