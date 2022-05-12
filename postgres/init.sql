-- 初始化

-- 主鍵
CREATE SEQUENCE user_seq
	START WITH 1
	INCREMENT BY 1
	NO MINVALUE
	MAXVALUE 2147483647
	CACHE 1;

-- 表
CREATE TABLE public.user
(
	id integer NOT NULL PRIMARY KEY DEFAULT nextval('user_seq'::regclass),
	name character varying(60) NOT NULL,
	mail character varying(32) NOT NULL,
	password character varying(64) NOT NULL,
	verified bool DEFAULT false NOT NULL,
	createtime timestamp(6) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatetime timestamp(6) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE public.user
    IS '帳號';

COMMENT ON COLUMN public.user.id
    IS '主鍵';

COMMENT ON COLUMN public.user.name
    IS '名稱';

COMMENT ON COLUMN public.user.mail
    IS '信箱';

COMMENT ON COLUMN public.user.password
    IS '密碼';

COMMENT ON COLUMN public.user.verified
    IS '是否驗證';

COMMENT ON COLUMN public.user.createtime
    IS '建立時間';

COMMENT ON COLUMN public.user.updatetime
    IS '更新時間';


-- 建立觸發器
CREATE OR REPLACE FUNCTION user_update_time() RETURNS TRIGGER AS 

$$
BEGIN
    NEW.updatetime = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER time_audit BEFORE UPDATE ON public.user
	FOR EACH ROW EXECUTE PROCEDURE user_update_time();
