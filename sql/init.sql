create table if not exists "users"
(
    user_id         integer not null 
                    constraint user_pk primary key,
    login           varchar(30),
    password        varchar(30),
    is_admin        boolean
);

create table if not exists rocket_flights
(
    flight_id       integer not null 
                    constraint flight_pk primary key,
    user_id         integer 
                    constraint flight_creator_user_id_fk references "users",
    moderator_id    integer 
                    constraint flight_moderator_user_id_fk references "users", 
    status          varchar(20),
    created_at      timestamp,
    formed_at       timestamp,
    confirmed_at    timestamp,
    flight_date     timestamp,
    payload         integer,
    price           float,
    title           varchar(100),
    site_number     integer
);

create table if not exists flight_requests
(
    request_id          integer not null 
                        constraint request_pk primary key,
    is_available        boolean,
    img_url             TEXT,
    title               varchar(100),
    load_capacity       float,
    description         TEXT,
    detailed_desc       TEXT,
    desired_price       float,
    flight_date_start   timestamp,
    flight_date_end     timestamp
);

create table if not exists flights_flight_requests
(
    flight_id       integer 
                    constraint flight_request_flight_flight_id_fk references rocket_flights,
    request_id      integer 
                    constraint request_flight_request_request_id_fk references flight_requests,
                primary key (flight_id, request_id)
);

alter table "users" owner to admin;
alter table rocket_flights owner to admin;
alter table flight_requests owner to admin;
alter table flights_flight_requests owner to admin;

INSERT INTO flight_requests(request_id, is_available, img_url, title, load_capacity,
 description, detailed_desc, desired_price, flight_date_start, flight_date_end)
VALUES (1, true, 'https://ntv-static.cdnvideo.ru/home/news/2023/20230205/sputn_io.jpg', '«Электро-Л» № 4', 1.8, 'Гидрометеорологическй космический аппарат',
    'Спутники «Электро-Л» создаются в рамках Федеральной космической программы России и входят' ||
		' в геостационарную гидрометеорологическую космическую систему «Электро» разработки НПО Лавочкина. Они предназначены' ||
		' для обеспечения оперативной и независимой гидрометеорологической информацией подразделений Федеральной службы по' ||
		' гидрометеорологии и мониторингу окружающей среды (Росгидромет) и других ведомств. ' ||
		'Сейчас в системе «Электро», функционирующей на околоземной орбите с 2011 года, используются по целевому назначению ' ||
		'два спутника — «Электро-Л» № 2 (запущен 11 декабря 2015 года) в точке стояния 14,5° западной долготы и «Электро-Л» №' ||
		' 3 (запущен 24 декабря 2019 года) в точке стояния 76° восточной долготы. Аппарату «Электро-Л» № 4 предстоит работать' ||
		' в точке стояния 165,8° восточной долготы. \n' ||
		'Уникальные особенности спутников «Электро-Л» позволяют получать независимые метеоданные с орбиты Земли каждые 15—30' ||
		' минут. Благодаря круглосуточной передаче с космических аппаратов высококачественных многоспектральных снимков повышается' ||
		' не только качество и оперативность прогнозов погоды, но и решаются глобальные вопросы мониторинга климата и изменений,' ||
		' выдаются штормовые и экстренные телеграммы при выявлении чрезвычайных ситуаций. \n' ||
		'Также, спутники ретранслируют сигналы от аварийных радиобуев международной спутниковой поисково-спасательной системы' ||
		' КОСПАС-САРСАТ. Это помогает поисково-спасательным службам эффективнее реагировать на сигналы бедствия для спасения человеческих жизней.', 
    16.5, TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'), TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'));
INSERT INTO flight_requests(request_id, is_available, img_url, title, load_capacity,
 description, detailed_desc, desired_price, flight_date_start, flight_date_end)
VALUES (2, true, 'https://upload.wikimedia.org/wikipedia/commons/thumb/b/b2/ISS-59_Progress_MS-11_approaches_the_ISS.jpg/1200px-ISS-59_Progress_MS-11_approaches_the_ISS.jpg',
	'«Прогресс МС-22»', 2.5, 'Грузовой корабль', 
	'«Прогресс МС-22» — космический транспортный грузовой корабль серии «Прогресс», запуск которого состоялся 6' ||
		' февраля 2023 года со стартового комплекса «Восток» (Площадка 31) космодрома «Байконур» по программе 83-й миссии снабжения ' ||
		'Международной космической станции[1]. Это 175-й полёт космического корабля серии «Прогресс».',
	26.5, 	TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'), TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'));
INSERT INTO flight_requests(request_id, is_available, img_url, title, load_capacity,
 description, detailed_desc, desired_price, flight_date_start, flight_date_end)
VALUES (3, true, 'https://news.cgtn.com/news/2023-02-26/Russia-s-replacement-Soyuz-spacecraft-arrives-at-space-station-1hJfEkcHidG/img/b5c3147c02e44da8941bf851fdebfcc7/b5c3147c02e44da8941bf851fdebfcc7-1280.png',
	'«Союз МС-23»', 7.2, 'Беспилотный корабль',
	'"Союз МС" ("МС" - "модернизированные системы") принадлежит к семейству космических' ||
		' кораблей "Союз" (первый пилотируемый запуск состоялся в 1967 году). Головным разработчиком и изготовителем' ||
		' корабля является Ракетно-космическая корпорация "Энергия" им. С. П. Королева (РКК "Энергия"; город Королев,' ||
		' Московская область). Эскизный проект "Союза МС", разработанный по заданию Федерального космического агентства' ||
		' (ныне госкорпорация "Роскосмос"), был одобрен на заседании научно-технического совета РКК "Энергия" в августе' ||
		' 2011 года. Корабль создан на базе предыдущей модификации "Союз ТМА-М" (запуски проводились в 2010-2016 годах)' ||
		' путем глубокой модернизации. "Союз МС" предназначен для доставки экипажей на МКС и возвращения их с орбиты' ||
		' обратно на Землю. Он выполняет роль корабля-спасателя в случаях вынужденной или аварийной эвакуации экипажа' ||
		' (при возникновении опасной ситуации на станции, заболевания или травмы космонавтов). Кроме того, "Союз МС"' ||
		' используется для доставки на МКС и возвращения с орбиты небольших грузов (научно-исследовательской аппаратуры' ||
		' и результатов экспериментов, личных вещей космонавтов и др.), а также для удаления со станции отходов в бытовом' ||
		' отсеке, который сгорает плотных слоях атмосферы при спуске корабля.',
	35.5, TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'), TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'));
INSERT INTO flight_requests(request_id, is_available, img_url, title, load_capacity,
 description, detailed_desc, desired_price, flight_date_start, flight_date_end)
VALUES (4, true, 'https://finobzor.ru/uploads/posts/2016-09/org_vrke626.jpg', '«Луч-5Х»', 37, 'Многофункциональная космическая система ретрансляции',
	'"Олимп-К", также обозначаемый как "Луч", является российским геостационарным спутником,' ||
		' созданным для Министерства обороны России и российского разведывательного управления ФСБ. Цели миссий не' ||
		' опубликованы. Согласно сообщению "Коммерсанта", спутник будет выполнять двойную роль: одна из них - радиотехническая' ||
		' разведка (SIGINT), а другая обеспечивает защищенную связь для правительственных нужд. Обозначение "Луч" указывает' ||
		' на роль ретранслятора данных. Следовательно, обозначение "Олимп-К" может относиться к полезной нагрузке' ||
		' SIGINT, в то время как обозначение "Луч" относится к полезной нагрузке ретрансляции данных. Другой источник' ||
		' сообщает, что спутник должен обеспечивать сигналы навигационной коррекции для системы ГЛОНАСС. Сообщалось также' ||
		' о бортовом лазерном устройстве связи.',
	65, TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'), TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'));
INSERT INTO flight_requests(request_id, is_available, img_url, title, load_capacity,
 description, detailed_desc, desired_price, flight_date_start, flight_date_end)
VALUES (5, true, 'https://avatars.dzeninfra.ru/get-zen_doc/9428044/pub_641e3138e540d5493c71189b_641e520f617db875869202c0/scale_1200',
	'«Барс-М №4»', 4, 'Электронно-оптический спутник',
	'Спутник "Барс-М" - это новый электронно-оптический спутник наблюдения за территорией, который' ||
		' заменит серию "Янтарь-1KFT" (Kometa) с возвратом пленки и отмененную серию "Барс".' ||
		' "Барс-М" является вторым воплощением проекта "Барс", который был начат в середине 1990-х годов для разработки преемника спутников' ||
		' наблюдения за территорией класса Komtea. Первоначальный проект Bars был остановлен в начале 2000-х годов. В' ||
		' 2007 году "ЦСКБ-Прогресс" заключило контракт на поставку "Барс-М", для которого, как сообщается, сервисный' ||
		' модуль на базе "Янтаря" был заменен новым усовершенствованным сервисным модулем.' ||
		' Спутники "Барс-М" оснащены электронно-оптической фотосистемой "Карат", разработанной и изготовленной Ленинградским' ||
		' оптико-механическим объединением (ЛОМО), и двойным лазерным высотомером для получения топографических изображений,' ||
		' стереоизображений, данных высотомера и изображений высокого разрешения с наземным разрешением около 1 метра.',
	17.4, TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'), TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'));

-- INSERT INTO flight_requests(request_id, is_available, img_url, title, load_capacity,
--  description, detailed_desc, desired_price, flight_date_start, flight_date_end)
-- VALUES (1, true, 'https://ntv-static.cdnvideo.ru/home/news/2023/20230205/sputn_io.jpg', '«Электро-Л» № 4', 1.8, 'Гидрометеорологическй космический аппарат',
--     'Спутники «Электро-Л» создаются в рамках Федеральной космической программы России и входят', 
--     16.5, TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'), TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'));
-- INSERT INTO flight_requests(request_id, is_available, img_url, title, load_capacity,
--  description, detailed_desc, desired_price, flight_date_start, flight_date_end)
-- VALUES (2, true, 'https://upload.wikimedia.org/wikipedia/commons/thumb/b/b2/ISS-59_Progress_MS-11_approaches_the_ISS.jpg/1200px-ISS-59_Progress_MS-11_approaches_the_ISS.jpg',
-- 	'«Прогресс МС-22»', 2.5, 'Грузовой корабль', 
-- 	'«Прогресс МС-22» — космический транспортный грузовой корабль серии «Прогресс», запуск которого состоялся 6' ||
-- 		' февраля 2023 года со стартового комплекса «Восток» (Площадка 31) космодрома «Байконур» по программе 83-й миссии снабжения ' ||
-- 		'Международной космической станции[1]. Это 175-й полёт космического корабля серии «Прогресс».',
-- 	26.5, 	TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'), TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'));
-- INSERT INTO flight_requests(request_id, is_available, img_url, title, load_capacity,
--  description, detailed_desc, desired_price, flight_date_start, flight_date_end)
-- VALUES (3, true, 'https://news.cgtn.com/news/2023-02-26/Russia-s-replacement-Soyuz-spacecraft-arrives-at-space-station-1hJfEkcHidG/img/b5c3147c02e44da8941bf851fdebfcc7/b5c3147c02e44da8941bf851fdebfcc7-1280.png',
-- 	'«Союз МС-23»', 7.2, 'Беспилотный корабль',
-- 	'"Союз МС" ("МС" - "модернизированные системы") принадлежит к семейству космических',
-- 	35.5, TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'), TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'));
-- INSERT INTO flight_requests(request_id, is_available, img_url, title, load_capacity,
--  description, detailed_desc, desired_price, flight_date_start, flight_date_end)
-- VALUES (4, true, 'https://finobzor.ru/uploads/posts/2016-09/org_vrke626.jpg', '«Луч-5Х»', 37, 'Многофункциональная космическая система ретрансляции',
-- 	'"Олимп-К", также обозначаемый как "Луч", является российским геостационарным спутником,' ||
-- 		' созданным для Министерства обороны России и российского разведывательного управления ФСБ. Цели миссий не',
-- 	65, TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'), TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'));
-- INSERT INTO flight_requests(request_id, is_available, img_url, title, load_capacity,
--  description, detailed_desc, desired_price, flight_date_start, flight_date_end)
-- VALUES (5, true, 'https://finobzor.ru/uploads/posts/2016-09/org_vrke626.jpg', '«Луч-5Х»', 37, 'Многофункциональная космическая система ретрансляции',
-- 	'"Олимп-К", также обозначаемый как "Луч", является российским геостационарным спутником,',
-- 	65, TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'), TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS')),
-- (6, true, 'https://avatars.dzeninfra.ru/get-zen_doc/9428044/pub_641e3138e540d5493c71189b_641e520f617db875869202c0/scale_1200',
-- 	'«Барс-М №4»', 4, 'Электронно-оптический спутник',
-- 	'Спутник "Барс-М" - это новый электронно-оптический спутник наблюдения за территорией, который',
-- 	17.4, TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'), TO_DATE('2023-02-05 12:12:52', 'YYYY-MM-DD HH24:MI:SS'));

