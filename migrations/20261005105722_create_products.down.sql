drop index if exists idx_products_slug on products (slug);

drop index if exists idx_products_category_id on products (category_id);

drop index if exists idx_products_created_by on products (created_by);

drop index if exists idx_products_status on products (status);

drop table if exists products;
