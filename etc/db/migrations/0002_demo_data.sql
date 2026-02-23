INSERT INTO stores (label, url)
VALUES
 ('3DFillies', 'https://3dfillies.com') -- 1
,('AnyCubic AU', 'https://anycubic.au') -- 2
,('CCDIY', 'https://ccdiy.com.au')  -- 3
,('CubicTech', 'https://cubictech.com.au') -- 4
,('eBay', 'https://ebay.com.au') -- 5
,('JayCar', 'https://jaycar.com.au') -- 6
,('Siddament', 'https://siddament.com.au') -- 7
,('Ali Express', 'https://aliexpress.com') -- 8
,('Prusament', 'https://prusament.com') -- 9
,('Amazon', 'https://amazon.com') -- 10
,('eSun', 'https://esun3dstore.com') -- 11
,('Jayo', 'https://jayo3d.com') -- 12
,('Creality AU', 'https://au.store.creality.com') --13
,('Sunlu', 'https://store.sunlu.com') -- 14
,('Elegoo', 'https://au.elegoo.com') -- 15
,('FlashForge', 'https://flashforge.com') -- 16
,('MatterHackers', 'https://matterhackers.com') -- 17
,('HobbyKing', 'https://hobbyking.com') -- 18
,('Creative Monkey', 'https://www.creativemonkey.com.au') -- 19
,('Ink Station', 'https://www.inkstation.com.au') -- 20
,('3D Central', 'https://3dcentral.com.au') -- 21
,('Polymaker', 'https://us.polymaker.com') -- 22
,('PhaserFPV', 'https://phaserfpv.com') -- 23
,('Filamentive', 'https://filamentive.com') -- 24
,('CopyMaster3D', 'https://copymaster3d.com') -- 25
,('Protopasta', 'https://protopasta.com') -- 26
,('3D Prima', 'https://3dprima.com') -- 27
,('3D Jake', 'https://3djake.com') -- 28
;

INSERT INTO Brands (label, active, store_id)
VALUES
 ('eSun', true, 11) -- 1
,('3DFillies', false, 1) -- 2
,('CCDIY', true, 3)  -- 3
,('Sunlu', true, 14) -- 4
,('AnyCubic', true, 2) -- 5
,('Creality', true, 13) -- 6
,('Elegoo', true, 15) -- 7
,('Siddament', true, 7) -- 8
,('Slic3d', true, NULL) -- 9
,('Jayo', true, 12) -- 10
;

INSERT INTO Locations (label, description, capacity, printable, tally)
VALUES
 ('Box 1', 'Top left most plastic storage box on the third shelf of the closet', 6, false, true) -- 1
,('Box 2', 'Top right most plastic storage box on the third shelf of the closet', 6, false, true) -- 2
,('Box 3', 'Top left most plastic storage box on the fourth shelf of the closet', 6, false, true) -- 3
,('Box 4', 'Top right most plastic storage box on the fourth shelf of the closet', 6, false, true) -- 4
,('Box 5', 'Bottom left most plastic storage box on the fourth shelf of the closet', 6, false, true) -- 5
,('Box 6', 'Bottom right most plastic storage box on the fourth shelf of the closet', 6, false, true) -- 6
,('Desk Dryer', 'Sovol SH02 on the desk beside the Kobra S1', 2, true, true) -- 7
,('ACE Pro', 'ACE Pro on Kobra S1', 4, true, true) -- 8
,('Order', 'Filament that needs to be ordered but has not yet been purchased', 0, false,true) -- 9
,('Shipping', 'Filament still being shipped after purchase', 0, false,true) -- 10
,('Limbo', 'No assigned storage space, probably on a desk or floor', 0, false,true) -- 11
,('Bin', 'Filament has been used or disposed of', 0, false, false) -- 12
;

INSERT INTO materials (label, class, description, special)
VALUES
 ('Basic PLA', 'PLA', 'Polylactic Acid', false) -- 1
,('ABS', 'ABS', 'Acrylonitrile Butadiene Styrene', true) -- 2
,('PETG', 'PETG', 'Polyethylene Terephthalate Glycol', false) -- 3
,('TPU 95A', 'TPU', 'Thermoplastic Polyurethane', true) -- 4
,('Hyper PLA', 'PLA', 'High-performance PLA blend', false) -- 5
,('High speed PLA', 'PLA', 'PLA optimized for high-speed printing', false) -- 6
,('PLA+', 'PLA', 'Enhanced PLA with improved properties', false) -- 7
,('PLA+ High speed', 'PLA', 'Enhanced PLA with improved properties', false) -- 8
,('Hyper PETG', 'PETG', 'High performance blend of PETG', false) -- 9
,('Wood PLA', 'PLA', 'PLA infused with wood fibers', true) -- 10
,('Silk PLA', 'PLA', 'PLA with a silky finish', false) -- 11
,('Matte PLA', 'PLA', 'Similar to PLA but the layer lines are less visible', false) -- 12
,('PA6 CF', 'Nylon', 'Carbon fibre reinforced nylon', true) -- 13
,('HS PLA+', 'PLA', 'High Speed PLA+', false) -- 14
;

INSERT INTO colors (label, hex_code, alias)
VALUES
 ('Black', '#000000', NULL) -- 1
,('Blue', '#2b3ac7', 'Klein Blue') -- 2
,('Brown', '#9f6f54', NULL) -- 3
,('Clear', '#FFF3', NULL) -- 4
,('Dark Blue', '#00263a', NULL) -- 5
,('Deep Green', '#07a590', NULL) -- 6
,('Firetruck Red', '#a50034', 'Aussie Toes, Fire Engine Red') -- 7
,('Gold', '#fed141', NULL) -- 8
,('Green', '#00a65a', NULL) -- 9
,('Grey', '#E0E0E0', NULL) -- 10
,('Interstellar Violet', '#5b618f', NULL) -- 11
,('Magenta', '#f14b70', NULL) -- 12
,('Teal', '#32cfb3', 'Mint Green') -- 13
,('Mint Green', '#a2e4b8', NULL) -- 14
,('Peach Pink', '#ffc196', NULL) -- 15
,('Green', '#00CC11', 'Rainbow') -- 16
,('Amber', '#E2BA1E', 'Rainbow') -- 17
,('Red', '#993300', 'Rainbow') -- 18
,('Red', '#ff1010', NULL) -- 19
,('Spring Leaf', '#89a84f', NULL) -- 20
,('Texture Grey', '#E6E6E6', NULL) -- 21
,('Translucent Red', '#fc1e1d', NULL) -- 22
,('Tropical Turquoise', '#009cbd', NULL) -- 23
,('White', '#FFFFFF', NULL) -- 24
,('Yellow', '#f8f852', NULL) -- 25
,('Orange', '#ff7f32', NULL) -- 26
,('Walnut', '#685e54', NULL) -- 27
,('Silver', '#858585', NULL) -- 28
,('Teal', '#01feb2', 'Soft Green') -- 29
,('Purple', '#aa5fde', NULL) -- 30
,('Royal Blue', '#38a2ec', NULL) -- 31
;

INSERT INTO spools (material_id, brand_id, location_id, store_id, weight, combined_weight, current_weight, price, empty, created_at, updated_at, deleted_at)
VALUES
 (7, 2, 12, 1, 1000.0, 0, 0, 0, true, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 1
,(7, 2, 8, 1, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2026-01-10 06:43:49.06431+00', NULL) -- 2
,(7, 2, 3, 1, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 3
,(7, 2, 3, 1, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 4
,(7, 2, 12, 1, 1000.0, 0, 0, 0, true, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 5
,(1, 5, 12, 2, 1000.0, 0, 0, 0, true, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 6
,(1, 5, 1, 2, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 7
,(1, 5, 1, 2, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 8
,(1, 5, 1, 2, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2026-01-11 03:26:15.821193+00', NULL) -- 9
,(1, 5, 12, 2, 1000.0, 0, 0, 0, true, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 10
,(1, 5, 2, 2, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 11
,(7, 4, 12, 6, 1000.0, 0, 0, 0, true, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 12
,(7, 4, 2, 6, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 13
,(7, 4, 2, 6, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 14
,(7, 4, 2, 6, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 15
,(7, 4, 8, 6, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 16
,(4, 7, 7, 6, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 17
,(6, 3, 12, 3, 1000.0, 0, 0, 0, true, '2025-12-13 06:02:33.865989+00', '2026-01-05 06:26:38.313367+00', NULL) -- 18
,(1, 3, 12, 3, 1000.0, 0, 0, 0, true, '2025-12-13 06:02:33.865989+00', '2026-01-08 06:00:55.016722+00', NULL) -- 19
,(2, 4, 3, 5, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 20
,(2, 4, 5, 5, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 21
,(7, 4, 12, 5, 1000.0, 0, 0, 0, true, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 22
,(7, 4, 12, 5, 1000.0, 0, 0, 0, true, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 23
,(7, 4, 12, 5, 1000.0, 0, 0, 0, true, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 24
,(7, 4, 12, 5, 1000.0, 0, 0, 0, true, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 25
,(7, 1, 4, 4, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 26
,(7, 1, 4, 4, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 27
,(7, 1, 4, 4, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 28
,(7, 1, 4, 4, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 29
,(5, 6, 12, 5, 1000.0, 0, 0, 0, true, '2025-12-13 06:02:33.865989+00', '2026-01-05 06:26:42.170084+00', NULL) -- 30
,(5, 6, 12, 5, 1000.0, 0, 0, 0, true, '2025-12-13 06:02:33.865989+00', '2026-01-05 06:26:45.904866+00', NULL) -- 31
,(5, 6, 12, 5, 1000.0, 0, 0, 0, true, '2025-12-13 06:02:33.865989+00', '2026-01-05 06:26:50.422653+00', NULL) -- 32
,(9, 6, 12, 5, 1000.0, 0, 0, 0, true, '2025-12-13 06:02:33.865989+00', '2026-01-09 13:18:54.19694+00', NULL) -- 33
,(5, 6, 3, 5, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 34
,(5, 6, 4, 5, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 35
,(1, 5, 12, 2, 1000.0, 0, 0, 0, true, '2025-12-13 06:02:33.865989+00', '2026-01-05 06:26:03.356671+00', NULL) -- 36
,(1, 5, 12, 2, 1000.0, 0, 0, 0, true, '2025-12-13 06:02:33.865989+00', '2026-01-05 06:25:57.274685+00', NULL) -- 37
,(3, 5, 5, 2, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 38
,(1, 5, 1, 2, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 39
,(1, 5, 1, 2, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 40
,(10, 4, 4, 2, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 41
,(11, 8, 1, 7, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 42
,(12, 8, 3, 7, 1000.0, 0, 0, 0, false, '2025-12-13 06:02:33.865989+00', '2025-12-13 06:02:33.865989+00', NULL) -- 43
,(3, 9, 5, 6, 1000.0, 0, 0, 0, false, '2025-12-22 21:19:49.652861+00', '2025-12-22 21:19:49.652861+00', NULL) -- 44
,(3, 9, 5, 6, 1000.0, 0, 0, 0, false, '2025-12-22 21:19:53.664251+00', '2025-12-22 21:19:53.664251+00', NULL) -- 45
,(3, 9, 8, 6, 1000.0, 0, 0, 0, false, '2025-12-22 21:19:54.922974+00', '2026-01-09 13:20:15.209033+00', NULL) -- 46
,(3, 9, 5, 6, 1000.0, 0, 0, 0, false, '2025-12-22 21:19:57.025819+00', '2025-12-22 21:19:57.025819+00', NULL) -- 47
,(1, 9, 12, 6, 1000.0, 0, 0, 0, true, '2025-12-22 21:20:58.264391+00', '2026-01-08 06:41:56.568233+00', NULL) -- 48
,(13, 4, 7, 6, 1000.0, 0, 0, 0, false, '2025-12-25 02:12:28.373172+00', '2025-12-25 02:12:28.373172+00', NULL) -- 49
,(1, 1, 8, 4, 1000.0, 0, 0, 0, false, '2025-12-25 05:41:01.77535+00', '2026-01-14 06:33:39.060886+00', NULL) -- 50
,(1, 1, 6, 4, 1000.0, 0, 0, 0, false, '2026-01-14 06:34:38.11096+00', '2025-12-25 05:41:01.77535+00', NULL) -- 51
,(1, 1, 6, 4, 1000.0, 0, 0, 0, false, '2026-01-14 06:34:44.44919+00', '2025-12-25 05:41:01.77535+00', NULL) -- 52
,(1, 1, 6, 4, 1000.0, 0, 0, 0, false, '2025-12-25 05:41:01.77535+00', '2026-01-14 06:35:03.13319+00', NULL) -- 53
,(1, 1, 6, 4, 1000.0, 0, 0, 0, false, '2025-12-25 05:41:01.77535+00', '2026-01-14 06:35:12.813413+00', NULL) -- 54
,(1, 5, 3, 8, 1000.0, 0, 0, 0, false, '2026-01-01 02:30:33.723068+00', '2026-01-01 02:30:33.723068+00', NULL) -- 55
,(1, 5, 2, 8, 1000.0, 0, 0, 0, false, '2026-01-01 02:34:02.54101+00', '2026-01-10 06:44:06.222988+00', NULL) -- 56
,(5, 6, 6, 5, 1000.0, 0, 0, 0, false, '2026-01-05 04:10:53.619585+00', '2026-01-05 04:10:53.619585+00', NULL) -- 57
,(9, 6, 5, 5, 1000.0, 0, 0, 0, false, '2026-01-05 04:11:02.275277+00', '2026-01-05 04:11:02.275277+00', NULL) -- 58
;

INSERT INTO spool_colors (spool_id, color_id)
VALUES

 (1, 1)
,(2, 9)
,(3, 12)
,(4, 3)
,(5, 24)
,(6, 11)
,(7, 15)
,(8, 20)
,(9, 23)
,(10, 21)
,(11, 24)
,(12, 1)
,(13, 7)
,(14, 2)
,(15, 8)
,(16, 25)
,(17, 6)
,(18, 1)
,(19, 18)
,(19, 16)
,(19, 17)
,(20, 1)
,(21, 24)
,(22, 1)
,(23, 1)
,(24, 13)
,(25, 13)
,(26, 14)
,(27, 5)
,(28, 14)
,(29, 5)
,(30, 1)
,(31, 1)
,(32, 1)
,(33, 24)
,(34, 10)
,(35, 19)
,(36, 21)
,(37, 21)
,(38, 4)
,(39, 26)
,(40, 22)
,(41, 27)
,(42, 28)
,(43, 29)
,(44, 1)
,(45, 1)
,(46, 24)
,(47, 24)
,(48, 1)
,(49, 1)
,(50, 1)
,(51, 1)
,(52, 1)
,(53, 1)
,(54, 1)
,(55, 21)
,(56, 21)
,(57, 31)
,(58, 1)
;

---- create above / drop below ----

TRUNCATE TABLE spool_colors CASCADE;
TRUNCATE TABLE spools CASCADE;
TRUNCATE TABLE colors CASCADE;
TRUNCATE TABLE materials CASCADE;
TRUNCATE TABLE locations CASCADE;
TRUNCATE TABLE brands CASCADE;
TRUNCATE TABLE stores CASCADE;
