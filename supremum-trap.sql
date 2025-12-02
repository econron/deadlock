-- ### supremumを共に取ってしまうパターン
-- （その結果、supremumをとってしまいその後insert intention lock waitinになり最後まで処理が走らないから排他ロックが解除されず、デッドロックになるパターンとか）
-- 一方、共有ロックも複数とって良いが、同じ範囲に対して共有ロックをとり、それぞれのトランザクションがupdateをしようとするとデッドロックになるパターンもある。
-- （insert,update,deleteは基本排他ロックを取ろうとする）
CREATE DATABASE gap_test;
USE gap_test;

DROP TABLE IF EXISTS gap_test;

CREATE TABLE gap_test (id INT PRIMARY KEY);
INSERT INTO gap_test VALUES (10);
-- これで 10 より後ろはすべて「Gap（Supremum）」になります

-- 実行手順

-- TX1
BEGIN;

--TX2
BEGIN;

--TX1
SELECT * FROM gap_test WHERE id = 20 FOR UPDATE;

--TX2
SELECT * FROM gap_test WHERE id = 30 FOR UPDATE;

--TX1
INSERT INTO gap_test VALUES (20);

--TX2
INSERT INTO gap_test VALUES (30);