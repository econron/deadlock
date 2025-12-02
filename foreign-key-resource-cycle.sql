-- ## 外部キー制約で子供と親を同時に更新した際に起きるパターン

-- また、外部キー制約でもデッドロックが起きるパターンがある。

-- 子を追加したついでに親も更新するパターン。

-- parent,childがいるとして、

-- TX1 -> child1 insert with parent 1 (parent 1にSロック)
-- TX2 -> child2 insert with parent 1（parent 1にSロック）
-- TX1 -> parent1 update (ここでXロックに昇格しようとするがTX2がSロックを持つ)
-- TX2 -> parent1 update（Xロックを要求するがすでにTX1が昇格待ちなのでTX2が終わらない。つまりSロックが手放されない）

CREATE DATABASE foreign_key_resource_cycle;
USE foreign_key_resource_cycle;

DROP TABLE IF EXISTS parent;
DROP TABLE IF EXISTS child;

-- 更新用のカラム (cnt) を追加
CREATE TABLE parent (id INT PRIMARY KEY, cnt INT DEFAULT 0);
CREATE TABLE child (id INT PRIMARY KEY, parent_id INT, FOREIGN KEY (parent_id) REFERENCES parent(id));

INSERT INTO parent VALUES (1, 0);
INSERT INTO child VALUES (1, 1);

-- TX1
BEGIN; -- 1
INSERT INTO child VALUES (2, 1); -- 3
UPDATE parent SET cnt = cnt + 1 WHERE id = 1; -- 5
COMMIT;

-- TX2
BEGIN; -- 2
INSERT INTO child VALUES (3, 1); -- 4 
UPDATE parent SET cnt = cnt + 1 WHERE id = 1; -- 6
COMMIT;