create table units
(
    crawl_date text, -- iso8601/rfc3339
    name       text,
    bed_min    real,
    bed_max    real,
    bath_min   real,
    bath_max   real,
    sqft_min   real,
    sqft_max   real,
    price_min  real,
    price_max  real
);