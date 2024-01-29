create table if not exists events (
  id bigserial primary key,
  event jsonb not null,
  user_id uuid not null
);