--
-- PostgreSQL database dump
--

\restrict Rr2t6b7WQknkGvXc2S1CNMWfRYabHLRoyNOjktaRvCu4TeBktota5A2syhYsApK

-- Dumped from database version 18.1
-- Dumped by pg_dump version 18.1

-- Started on 2026-01-04 23:04:28

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 6 (class 2615 OID 18992)
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

-- *not* creating schema, since initdb creates it


ALTER SCHEMA public OWNER TO postgres;

--
-- TOC entry 5162 (class 0 OID 0)
-- Dependencies: 6
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON SCHEMA public IS '';


--
-- TOC entry 2 (class 3079 OID 18993)
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- TOC entry 5164 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 227 (class 1259 OID 19067)
-- Name: categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categories (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.categories OWNER TO postgres;

--
-- TOC entry 226 (class 1259 OID 19066)
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.categories_id_seq OWNER TO postgres;

--
-- TOC entry 5165 (class 0 OID 0)
-- Dependencies: 226
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- TOC entry 233 (class 1259 OID 19122)
-- Name: items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.items (
    id integer NOT NULL,
    sku character varying(50) NOT NULL,
    name character varying(150) NOT NULL,
    category_id integer NOT NULL,
    rack_id integer NOT NULL,
    stock integer DEFAULT 0 NOT NULL,
    minimum_stock integer DEFAULT 5 NOT NULL,
    price numeric(15,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT items_price_check CHECK ((price >= (0)::numeric))
);


ALTER TABLE public.items OWNER TO postgres;

--
-- TOC entry 232 (class 1259 OID 19121)
-- Name: items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.items_id_seq OWNER TO postgres;

--
-- TOC entry 5166 (class 0 OID 0)
-- Dependencies: 232
-- Name: items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.items_id_seq OWNED BY public.items.id;


--
-- TOC entry 231 (class 1259 OID 19099)
-- Name: racks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.racks (
    id integer NOT NULL,
    warehouse_id integer NOT NULL,
    code character varying(50) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.racks OWNER TO postgres;

--
-- TOC entry 230 (class 1259 OID 19098)
-- Name: racks_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.racks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.racks_id_seq OWNER TO postgres;

--
-- TOC entry 5167 (class 0 OID 0)
-- Dependencies: 230
-- Name: racks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.racks_id_seq OWNED BY public.racks.id;


--
-- TOC entry 221 (class 1259 OID 19005)
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.roles (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.roles OWNER TO postgres;

--
-- TOC entry 220 (class 1259 OID 19004)
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.roles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.roles_id_seq OWNER TO postgres;

--
-- TOC entry 5168 (class 0 OID 0)
-- Dependencies: 220
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- TOC entry 237 (class 1259 OID 19174)
-- Name: sale_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sale_items (
    id integer NOT NULL,
    sale_id integer NOT NULL,
    item_id integer NOT NULL,
    quantity integer NOT NULL,
    price_at_sale numeric(15,2) NOT NULL,
    subtotal numeric(15,2) NOT NULL,
    CONSTRAINT sale_items_price_at_sale_check CHECK ((price_at_sale >= (0)::numeric)),
    CONSTRAINT sale_items_quantity_check CHECK ((quantity > 0)),
    CONSTRAINT sale_items_subtotal_check CHECK ((subtotal >= (0)::numeric))
);


ALTER TABLE public.sale_items OWNER TO postgres;

--
-- TOC entry 236 (class 1259 OID 19173)
-- Name: sale_items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.sale_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.sale_items_id_seq OWNER TO postgres;

--
-- TOC entry 5169 (class 0 OID 0)
-- Dependencies: 236
-- Name: sale_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.sale_items_id_seq OWNED BY public.sale_items.id;


--
-- TOC entry 235 (class 1259 OID 19156)
-- Name: sales; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sales (
    id integer NOT NULL,
    user_id integer NOT NULL,
    total_amount numeric(15,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT sales_total_amount_check CHECK ((total_amount >= (0)::numeric))
);


ALTER TABLE public.sales OWNER TO postgres;

--
-- TOC entry 234 (class 1259 OID 19155)
-- Name: sales_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.sales_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.sales_id_seq OWNER TO postgres;

--
-- TOC entry 5170 (class 0 OID 0)
-- Dependencies: 234
-- Name: sales_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.sales_id_seq OWNED BY public.sales.id;


--
-- TOC entry 225 (class 1259 OID 19046)
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sessions (
    id integer NOT NULL,
    user_id integer NOT NULL,
    token uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    expired_at timestamp with time zone NOT NULL,
    revoked_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.sessions OWNER TO postgres;

--
-- TOC entry 224 (class 1259 OID 19045)
-- Name: sessions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.sessions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.sessions_id_seq OWNER TO postgres;

--
-- TOC entry 5171 (class 0 OID 0)
-- Dependencies: 224
-- Name: sessions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.sessions_id_seq OWNED BY public.sessions.id;


--
-- TOC entry 223 (class 1259 OID 19018)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    email character varying(100) NOT NULL,
    password_hash text NOT NULL,
    role_id integer NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 222 (class 1259 OID 19017)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- TOC entry 5172 (class 0 OID 0)
-- Dependencies: 222
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 229 (class 1259 OID 19084)
-- Name: warehouses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.warehouses (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    location text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.warehouses OWNER TO postgres;

--
-- TOC entry 228 (class 1259 OID 19083)
-- Name: warehouses_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.warehouses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.warehouses_id_seq OWNER TO postgres;

--
-- TOC entry 5173 (class 0 OID 0)
-- Dependencies: 228
-- Name: warehouses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.warehouses_id_seq OWNED BY public.warehouses.id;


--
-- TOC entry 4916 (class 2604 OID 19070)
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- TOC entry 4925 (class 2604 OID 19125)
-- Name: items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items ALTER COLUMN id SET DEFAULT nextval('public.items_id_seq'::regclass);


--
-- TOC entry 4922 (class 2604 OID 19102)
-- Name: racks id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.racks ALTER COLUMN id SET DEFAULT nextval('public.racks_id_seq'::regclass);


--
-- TOC entry 4907 (class 2604 OID 19008)
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- TOC entry 4933 (class 2604 OID 19177)
-- Name: sale_items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_items ALTER COLUMN id SET DEFAULT nextval('public.sale_items_id_seq'::regclass);


--
-- TOC entry 4930 (class 2604 OID 19159)
-- Name: sales id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales ALTER COLUMN id SET DEFAULT nextval('public.sales_id_seq'::regclass);


--
-- TOC entry 4913 (class 2604 OID 19049)
-- Name: sessions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions ALTER COLUMN id SET DEFAULT nextval('public.sessions_id_seq'::regclass);


--
-- TOC entry 4909 (class 2604 OID 19021)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 4919 (class 2604 OID 19087)
-- Name: warehouses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.warehouses ALTER COLUMN id SET DEFAULT nextval('public.warehouses_id_seq'::regclass);


--
-- TOC entry 5146 (class 0 OID 19067)
-- Dependencies: 227
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.categories (id, name, description, created_at, updated_at) FROM stdin;
1	Electronics	Electronic devices and accessories	2026-01-01 21:51:26.963183+07	2026-01-01 21:51:26.963183+07
2	Furniture	Office and home furniture	2026-01-01 21:51:26.963183+07	2026-01-01 21:51:26.963183+07
3	Office Supplies	Stationery and office equipment	2026-01-01 21:51:26.963183+07	2026-01-01 21:51:26.963183+07
4	Hardware	Tools and hardware equipment	2026-01-01 21:51:26.963183+07	2026-01-01 21:51:26.963183+07
5	Clothing	Work uniforms and clothing	2026-01-01 21:51:26.963183+07	2026-01-01 21:51:26.963183+07
6	Testing	Testing description update lagi 1	2026-01-02 00:22:55.702301+07	2026-01-02 00:29:10.905241+07
\.


--
-- TOC entry 5152 (class 0 OID 19122)
-- Dependencies: 233
-- Data for Name: items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.items (id, sku, name, category_id, rack_id, stock, minimum_stock, price, created_at, updated_at) FROM stdin;
4	FRN-001	Office Chair Ergonomic	2	2	30	5	1500000.00	2026-01-01 21:51:42.12791+07	2026-01-01 21:51:42.12791+07
5	FRN-002	Office Desk Standing	2	2	2	5	2500000.00	2026-01-01 21:51:42.12791+07	2026-01-01 21:51:42.12791+07
6	OFS-001	Printer HP LaserJet	3	3	25	5	3200000.00	2026-01-01 21:51:42.12791+07	2026-01-01 21:51:42.12791+07
8	OFS-003	Pen Set	3	3	4	10	25000.00	2026-01-01 21:51:42.12791+07	2026-01-01 21:51:42.12791+07
9	HRD-001	Hammer Tool Set	4	4	80	15	350000.00	2026-01-01 21:51:42.12791+07	2026-01-01 21:51:42.12791+07
10	CLO-001	Work Uniform Shirt	5	5	100	20	125000.00	2026-01-01 21:51:42.12791+07	2026-01-01 21:51:42.12791+07
11	NEW-001-CREATE	Laptop Asus ROG	1	1	25	5	15000000.00	2026-01-02 10:36:17.783293+07	2026-01-02 10:36:17.783293+07
3	ELC-003	Keyboard Mechanical	1	1	3	5	750000.00	2026-01-01 21:51:42.12791+07	2026-01-02 17:25:20.609774+07
1	ELC-001	Laptop Dell Inspiron 15	1	1	47	10	8500000.00	2026-01-01 21:51:42.12791+07	2026-01-04 21:24:42.765658+07
2	ELC-002	Mouse Wireless Logitech	1	1	142	20	250000.00	2026-01-01 21:51:42.12791+07	2026-01-04 21:24:42.765658+07
7	OFS-002	Paper A4 Sidu (1 Rim)	3	3	190	50	45000.00	2026-01-01 21:51:42.12791+07	2026-01-04 21:24:42.765658+07
\.


--
-- TOC entry 5150 (class 0 OID 19099)
-- Dependencies: 231
-- Data for Name: racks; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.racks (id, warehouse_id, code, description, created_at, updated_at) FROM stdin;
1	1	A-01	Electronics section row A	2026-01-01 21:51:26.963183+07	2026-01-01 21:51:26.963183+07
2	1	B-01	Furniture section row B	2026-01-01 21:51:26.963183+07	2026-01-01 21:51:26.963183+07
3	2	C-01	Office supplies section	2026-01-01 21:51:26.963183+07	2026-01-01 21:51:26.963183+07
4	3	D-01	Hardware section	2026-01-01 21:51:26.963183+07	2026-01-01 21:51:26.963183+07
5	4	E-01	Clothing section	2026-01-01 21:51:26.963183+07	2026-01-01 21:51:26.963183+07
6	5	A-01-NEW-UPDATE	Rak baru untuk elektronik test Update	2026-01-02 00:32:51.204667+07	2026-01-02 00:34:04.732614+07
\.


--
-- TOC entry 5140 (class 0 OID 19005)
-- Dependencies: 221
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.roles (id, name, created_at) FROM stdin;
1	super_admin	2026-01-01 21:51:05.861677+07
2	admin	2026-01-01 21:51:05.861677+07
3	staff	2026-01-01 21:51:05.861677+07
\.


--
-- TOC entry 5156 (class 0 OID 19174)
-- Dependencies: 237
-- Data for Name: sale_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sale_items (id, sale_id, item_id, quantity, price_at_sale, subtotal) FROM stdin;
1	1	2	5	250000.00	1250000.00
6	2	1	1	8500000.00	8500000.00
10	5	1	2	8500000.00	17000000.00
11	5	2	3	250000.00	750000.00
12	5	7	10	45000.00	450000.00
\.


--
-- TOC entry 5154 (class 0 OID 19156)
-- Dependencies: 235
-- Data for Name: sales; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sales (id, user_id, total_amount, created_at, updated_at, deleted_at) FROM stdin;
1	1	1250000.00	2026-01-02 16:47:53.14903+07	2026-01-02 17:16:23.675592+07	\N
2	1	8500000.00	2026-01-02 16:49:54.862416+07	2026-01-02 17:20:15.574819+07	\N
3	1	17000000.00	2026-01-02 17:20:47.753107+07	2026-01-02 17:23:33.979262+07	2026-01-02 17:23:46.834756+07
4	1	750000.00	2026-01-02 17:24:17.745816+07	2026-01-02 17:24:17.745816+07	2026-01-02 17:25:20.609774+07
5	1	18200000.00	2026-01-04 21:24:42.765658+07	2026-01-04 21:24:42.765658+07	\N
\.


--
-- TOC entry 5144 (class 0 OID 19046)
-- Dependencies: 225
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sessions (id, user_id, token, expired_at, revoked_at, created_at) FROM stdin;
1	1	40d5e4ab-38a8-42a9-875d-ec03fe16fd66	2026-01-02 22:32:22.018634+07	2026-01-01 22:37:40.010843+07	2026-01-01 22:32:22.01995+07
2	1	34dc1b1d-6959-492c-a0dd-f0532b1b0b08	2026-01-02 22:38:25.832542+07	2026-01-01 22:38:34.529679+07	2026-01-01 22:38:25.832917+07
3	1	a017e44d-56c8-401d-ab4b-a63545c45f7b	2026-01-02 23:24:58.640186+07	\N	2026-01-01 23:24:58.640918+07
4	1	d1835bd4-5b83-4275-a743-180cdf488740	2026-01-02 23:26:49.89041+07	2026-01-01 23:37:05.332976+07	2026-01-01 23:26:49.891419+07
5	1	730e6d13-774e-4be2-806c-87a09a7e4ff5	2026-01-02 23:27:38.797135+07	2026-01-01 23:37:10.667969+07	2026-01-01 23:27:38.797637+07
6	1	79eea459-6e21-425d-9b79-8e2d97ba2af9	2026-01-02 23:37:17.842055+07	\N	2026-01-01 23:37:17.842703+07
7	1	271b1f60-b111-4705-99a8-76e193b65000	2026-01-02 23:56:54.500535+07	\N	2026-01-01 23:56:54.50189+07
8	1	67ac6fc2-27f9-4f79-87e4-3553b88c5a8a	2026-01-04 15:07:11.143818+07	\N	2026-01-03 15:07:11.145995+07
9	1	27ea3b40-8fe8-4e2e-87ae-d4a242b2799e	2026-01-05 20:31:48.663384+07	2026-01-04 20:34:51.515066+07	2026-01-04 20:31:48.664972+07
10	1	003d2193-79b1-41c7-84de-a68ee202e8fe	2026-01-05 20:32:59.374497+07	2026-01-04 20:37:58.851032+07	2026-01-04 20:32:59.375624+07
11	1	c8ddfeff-49e2-42a0-9264-24a4f46e8531	2026-01-05 20:38:01.528367+07	2026-01-04 20:59:43.757004+07	2026-01-04 20:38:01.528645+07
12	1	3a808382-6a59-45be-9833-872cbcbd76f7	2026-01-05 20:59:53.182985+07	\N	2026-01-04 20:59:53.1843+07
13	1	0bebf396-2418-4eec-8bda-87af6b3b4318	2026-01-05 21:00:44.254657+07	\N	2026-01-04 21:00:44.255083+07
14	1	20e6d5ca-0815-4ce9-9f99-720fb6bc98b7	2026-01-05 21:05:03.959141+07	2026-01-04 21:11:30.015598+07	2026-01-04 21:05:03.960143+07
15	1	f3382e77-abbe-4c38-bcd4-f43b9b200888	2026-01-05 21:11:22.320386+07	2026-01-04 21:15:52.401211+07	2026-01-04 21:11:22.321276+07
16	1	6272f7a8-cb1e-403a-b3f3-dc4153965cb7	2026-01-05 21:15:57.023219+07	\N	2026-01-04 21:15:57.024243+07
\.


--
-- TOC entry 5142 (class 0 OID 19018)
-- Dependencies: 223
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name, email, password_hash, role_id, is_active, created_at, updated_at) FROM stdin;
1	Super Admin User	superadmin@inventory.com	$2a$10$GOArleP8YiH7jncpDEWNQORUlH25w9v28MB9Bylfe2pjZlYPZUGYm	1	t	2026-01-01 22:32:12.555356+07	2026-01-01 22:32:12.555356+07
2	Admin User	admin@inventory.com	$2a$10$GOArleP8YiH7jncpDEWNQORUlH25w9v28MB9Bylfe2pjZlYPZUGYm	2	t	2026-01-01 22:32:12.555356+07	2026-01-01 22:32:12.555356+07
3	Staff User	staff@inventory.com	$2a$10$GOArleP8YiH7jncpDEWNQORUlH25w9v28MB9Bylfe2pjZlYPZUGYm	3	t	2026-01-01 22:32:12.555356+07	2026-01-01 22:32:12.555356+07
4	Updated Staff Name	staff2.updated@inventory.com	$2a$10$.zvSl/puE7Om9FwCWoZYTOHXOpV0lt78acwp1aves6pFCDc3sgG7K	2	f	2026-01-02 17:31:04.145346+07	2026-01-02 17:35:17.176736+07
6	Test Deleted	deleted@inventory.com	$2a$10$J/1K.1RAMFE.4ZdduFOPje9aemVBIkkzhFj.7bKscIZ5/uCllbyUK	3	t	2026-01-04 21:23:49.713559+07	2026-01-04 21:23:49.713559+07
\.


--
-- TOC entry 5148 (class 0 OID 19084)
-- Dependencies: 229
-- Data for Name: warehouses; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.warehouses (id, name, location, created_at, updated_at) FROM stdin;
1	Main Warehouse	Jl. Industri No. 1, Jakarta Utara	2026-01-01 21:51:26.963183+07	2026-01-01 21:51:26.963183+07
2	South Warehouse	Jl. Raya Selatan No. 45, Jakarta Selatan	2026-01-01 21:51:26.963183+07	2026-01-01 21:51:26.963183+07
3	East Warehouse	Jl. Timur Raya No. 88, Bekasi	2026-01-01 21:51:26.963183+07	2026-01-01 21:51:26.963183+07
4	West Warehouse	Jl. Barat No. 12, Tangerang	2026-01-01 21:51:26.963183+07	2026-01-01 21:51:26.963183+07
5	Central Warehouse	Jl. Pusat No. 77, Jakarta Pusat	2026-01-01 21:51:26.963183+07	2026-01-01 21:51:26.963183+07
6	Test Warehouse	Waru, Sidoarjo, Jawa Timur	2026-01-02 16:36:51.676424+07	2026-01-02 16:36:51.676424+07
\.


--
-- TOC entry 5174 (class 0 OID 0)
-- Dependencies: 226
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.categories_id_seq', 9, true);


--
-- TOC entry 5175 (class 0 OID 0)
-- Dependencies: 232
-- Name: items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.items_id_seq', 13, true);


--
-- TOC entry 5176 (class 0 OID 0)
-- Dependencies: 230
-- Name: racks_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.racks_id_seq', 7, true);


--
-- TOC entry 5177 (class 0 OID 0)
-- Dependencies: 220
-- Name: roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.roles_id_seq', 3, true);


--
-- TOC entry 5178 (class 0 OID 0)
-- Dependencies: 236
-- Name: sale_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sale_items_id_seq', 12, true);


--
-- TOC entry 5179 (class 0 OID 0)
-- Dependencies: 234
-- Name: sales_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sales_id_seq', 5, true);


--
-- TOC entry 5180 (class 0 OID 0)
-- Dependencies: 224
-- Name: sessions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sessions_id_seq', 16, true);


--
-- TOC entry 5181 (class 0 OID 0)
-- Dependencies: 222
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 6, true);


--
-- TOC entry 5182 (class 0 OID 0)
-- Dependencies: 228
-- Name: warehouses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.warehouses_id_seq', 8, true);


--
-- TOC entry 4957 (class 2606 OID 19082)
-- Name: categories categories_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_name_key UNIQUE (name);


--
-- TOC entry 4959 (class 2606 OID 19080)
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- TOC entry 4972 (class 2606 OID 19142)
-- Name: items items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_pkey PRIMARY KEY (id);


--
-- TOC entry 4974 (class 2606 OID 19144)
-- Name: items items_sku_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_sku_key UNIQUE (sku);


--
-- TOC entry 4964 (class 2606 OID 19113)
-- Name: racks racks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.racks
    ADD CONSTRAINT racks_pkey PRIMARY KEY (id);


--
-- TOC entry 4940 (class 2606 OID 19016)
-- Name: roles roles_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);


--
-- TOC entry 4942 (class 2606 OID 19014)
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- TOC entry 4983 (class 2606 OID 19188)
-- Name: sale_items sale_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT sale_items_pkey PRIMARY KEY (id);


--
-- TOC entry 4979 (class 2606 OID 19167)
-- Name: sales sales_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT sales_pkey PRIMARY KEY (id);


--
-- TOC entry 4953 (class 2606 OID 19058)
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- TOC entry 4955 (class 2606 OID 19060)
-- Name: sessions sessions_token_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_token_key UNIQUE (token);


--
-- TOC entry 4966 (class 2606 OID 19115)
-- Name: racks uq_rack_code_per_warehouse; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.racks
    ADD CONSTRAINT uq_rack_code_per_warehouse UNIQUE (warehouse_id, code);


--
-- TOC entry 4946 (class 2606 OID 19038)
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- TOC entry 4948 (class 2606 OID 19036)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 4961 (class 2606 OID 19097)
-- Name: warehouses warehouses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouses_pkey PRIMARY KEY (id);


--
-- TOC entry 4967 (class 1259 OID 19204)
-- Name: idx_items_category_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_items_category_id ON public.items USING btree (category_id);


--
-- TOC entry 4968 (class 1259 OID 19205)
-- Name: idx_items_rack_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_items_rack_id ON public.items USING btree (rack_id);


--
-- TOC entry 4969 (class 1259 OID 19207)
-- Name: idx_items_sku; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_items_sku ON public.items USING btree (sku);


--
-- TOC entry 4970 (class 1259 OID 19206)
-- Name: idx_items_stock; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_items_stock ON public.items USING btree (stock);


--
-- TOC entry 4962 (class 1259 OID 19208)
-- Name: idx_racks_warehouse_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_racks_warehouse_id ON public.racks USING btree (warehouse_id);


--
-- TOC entry 4980 (class 1259 OID 19212)
-- Name: idx_sale_items_item_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sale_items_item_id ON public.sale_items USING btree (item_id);


--
-- TOC entry 4981 (class 1259 OID 19211)
-- Name: idx_sale_items_sale_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sale_items_sale_id ON public.sale_items USING btree (sale_id);


--
-- TOC entry 4975 (class 1259 OID 19210)
-- Name: idx_sales_created_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sales_created_at ON public.sales USING btree (created_at);


--
-- TOC entry 4976 (class 1259 OID 19240)
-- Name: idx_sales_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sales_deleted_at ON public.sales USING btree (deleted_at) WHERE (deleted_at IS NOT NULL);


--
-- TOC entry 4977 (class 1259 OID 19209)
-- Name: idx_sales_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sales_user_id ON public.sales USING btree (user_id);


--
-- TOC entry 4949 (class 1259 OID 19203)
-- Name: idx_sessions_expired_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sessions_expired_at ON public.sessions USING btree (expired_at);


--
-- TOC entry 4950 (class 1259 OID 19202)
-- Name: idx_sessions_token; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sessions_token ON public.sessions USING btree (token);


--
-- TOC entry 4951 (class 1259 OID 19201)
-- Name: idx_sessions_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sessions_user_id ON public.sessions USING btree (user_id);


--
-- TOC entry 4943 (class 1259 OID 19200)
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- TOC entry 4944 (class 1259 OID 19199)
-- Name: idx_users_role_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_role_id ON public.users USING btree (role_id);


--
-- TOC entry 4987 (class 2606 OID 19145)
-- Name: items fk_items_category; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT fk_items_category FOREIGN KEY (category_id) REFERENCES public.categories(id);


--
-- TOC entry 4988 (class 2606 OID 19150)
-- Name: items fk_items_rack; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT fk_items_rack FOREIGN KEY (rack_id) REFERENCES public.racks(id);


--
-- TOC entry 4986 (class 2606 OID 19116)
-- Name: racks fk_racks_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.racks
    ADD CONSTRAINT fk_racks_warehouse FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id) ON DELETE CASCADE;


--
-- TOC entry 4990 (class 2606 OID 19194)
-- Name: sale_items fk_sale_items_item; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT fk_sale_items_item FOREIGN KEY (item_id) REFERENCES public.items(id);


--
-- TOC entry 4991 (class 2606 OID 19189)
-- Name: sale_items fk_sale_items_sale; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT fk_sale_items_sale FOREIGN KEY (sale_id) REFERENCES public.sales(id) ON DELETE CASCADE;


--
-- TOC entry 4989 (class 2606 OID 19168)
-- Name: sales fk_sales_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT fk_sales_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 4985 (class 2606 OID 19061)
-- Name: sessions fk_sessions_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 4984 (class 2606 OID 19039)
-- Name: users fk_users_role; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES public.roles(id) ON UPDATE CASCADE;


--
-- TOC entry 5163 (class 0 OID 0)
-- Dependencies: 6
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;


-- Completed on 2026-01-04 23:04:28

--
-- PostgreSQL database dump complete
--

\unrestrict Rr2t6b7WQknkGvXc2S1CNMWfRYabHLRoyNOjktaRvCu4TeBktota5A2syhYsApK

