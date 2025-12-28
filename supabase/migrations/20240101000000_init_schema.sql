-- Create profiles table
create table public.profiles (
  id uuid references auth.users not null primary key,
  username text,
  avatar_url text,
  bio text,
  phone text,
  preferences jsonb default '{}'::jsonb,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

-- Enable RLS for profiles
alter table public.profiles enable row level security;

create policy "Users can view own profile" on public.profiles
  for select using (auth.uid() = id);

create policy "Users can update own profile" on public.profiles
  for update using (auth.uid() = id);

create policy "Users can insert own profile" on public.profiles
  for insert with check (auth.uid() = id);

-- Handle new user trigger
create or replace function public.handle_new_user()
returns trigger as $$
begin
  insert into public.profiles (id, username, avatar_url)
  values (new.id, new.raw_user_meta_data->>'username', new.raw_user_meta_data->>'avatar_url');
  return new;
end;
$$ language plpgsql security definer;

-- Trigger logic is specific to Supabase Auth, assumes auth.users exists
drop trigger if exists on_auth_user_created on auth.users;
create trigger on_auth_user_created
  after insert on auth.users
  for each row execute procedure public.handle_new_user();

-- Create books table
create table public.books (
  id uuid default gen_random_uuid() primary key,
  user_id uuid references auth.users not null,
  title text not null,
  author text,
  cover_url text,
  file_url text not null,
  format text,
  size bigint,
  total_pages int,
  publisher text,
  description text,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

-- Enable RLS for books
alter table public.books enable row level security;

create policy "Users can view own books" on public.books
  for select using (auth.uid() = user_id);

create policy "Users can insert own books" on public.books
  for insert with check (auth.uid() = user_id);

create policy "Users can update own books" on public.books
  for update using (auth.uid() = user_id);

create policy "Users can delete own books" on public.books
  for delete using (auth.uid() = user_id);

-- Create reading_progress table
create table public.reading_progress (
  id uuid default gen_random_uuid() primary key,
  user_id uuid references auth.users not null,
  book_id uuid references public.books(id) on delete cascade not null,
  progress float default 0,
  current_loc text,
  status text check (status in ('reading', 'finished')),
  updated_at timestamptz default now(),
  unique(user_id, book_id)
);

-- Enable RLS for reading_progress
alter table public.reading_progress enable row level security;

create policy "Users can view own progress" on public.reading_progress
  for select using (auth.uid() = user_id);

create policy "Users can insert own progress" on public.reading_progress
  for insert with check (auth.uid() = user_id);

create policy "Users can update own progress" on public.reading_progress
  for update using (auth.uid() = user_id);

-- Create reading_sessions table
create table public.reading_sessions (
  id uuid default gen_random_uuid() primary key,
  user_id uuid references auth.users not null,
  book_id uuid references public.books(id) on delete cascade not null,
  start_time timestamptz not null,
  end_time timestamptz not null,
  duration int not null,
  created_at timestamptz default now()
);

-- Enable RLS for reading_sessions
alter table public.reading_sessions enable row level security;

create policy "Users can view own sessions" on public.reading_sessions
  for select using (auth.uid() = user_id);

create policy "Users can insert own sessions" on public.reading_sessions
  for insert with check (auth.uid() = user_id);

-- Storage Policies (Requires 'storage' schema enabled)
-- Note: You must create 'book-files' and 'book-covers' buckets in Supabase dashboard first

-- Enable RLS for objects in storage.objects
-- NOTE: RLS is usually enabled by default on storage.objects. 
-- If you get "must be owner of table objects", skip this line.
-- alter table storage.objects enable row level security;

drop policy if exists "Users can upload own book files" on storage.objects;
create policy "Users can upload own book files"
on storage.objects for insert
with check (
  bucket_id = 'book-files' AND
  auth.uid()::text = (storage.foldername(name))[1]
);

drop policy if exists "Users can view own book files" on storage.objects;
create policy "Users can view own book files"
on storage.objects for select
using (
  bucket_id = 'book-files' AND
  auth.uid()::text = (storage.foldername(name))[1]
);

drop policy if exists "Users can delete own book files" on storage.objects;
create policy "Users can delete own book files"
on storage.objects for delete
using (
  bucket_id = 'book-files' AND
  auth.uid()::text = (storage.foldername(name))[1]
);

-- Function: Get Reading Stats (Past 30 Days)
create or replace function get_reading_stats(days int default 30)
returns json
language plpgsql
security definer
as $$
declare
  result json;
begin
  select json_build_object(
    'total_duration', coalesce(sum(duration), 0),
    'books_finished', (select count(*) from reading_progress where user_id = auth.uid() and status = 'finished'),
    'daily_stats', (
      select json_agg(daily)
      from (
        select
          d.day::date as date,
          coalesce(sum(s.duration), 0) as duration
        from
          generate_series(now() - (days || ' days')::interval, now(), '1 day'::interval) as d(day)
          left join reading_sessions s on s.user_id = auth.uid() and s.created_at::date = d.day::date
        group by d.day
        order by d.day
      ) daily
    )
  ) into result
  from reading_sessions
  where user_id = auth.uid();
  
  return result;
end;
$$;
