-- +goose Up
-- Insert mock data
INSERT INTO
    public.items (name, brand_id, description, sex, category_id, price, discount, outer_link, created_at)
VALUES
    ('Кардиган DEVIL LOGO', 3, 'Кардиган в исполнении российского streetwear бренда Daze', '1', '6', '28000', '0', 'https://nuw.store/catalog/svitery_i_kardigany/kardigan_devil_logo_1-color-raznocvetnyj/', CURRENT_TIMESTAMP);

INSERT INTO
    public.items (name, brand_id, description, sex, category_id, price, discount, outer_link, created_at)
VALUES
    ('Худи WWW.S DOUBLE SIDED', 1, 'Худи в исполнении Gleb Kostin Solutions', '1', '7', '14780', '10', 'https://nuw.store/catalog/tolstovki/khudi_www_s_double_sided-color-chernyj/', CURRENT_TIMESTAMP);

INSERT INTO
    public.items (name, brand_id, description, sex, category_id, price, discount, outer_link, created_at)
VALUES
    ('Кроссовки MXR TECH', 2, 'Модель MXR отражает наше стремление распространить опыт ретро-бега Mizuno 2000-х годов в качестве дополнения к всемирно известной модели Wave Rider 10.', '1', '12', '15990', '0', 'https://nuw.store/catalog/krossovki/krossovki_mxr_tech-color-raznocvetnyj/', CURRENT_TIMESTAMP);

