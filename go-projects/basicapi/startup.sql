-- Create table with some dummy data
CREATE TABLE IF NOT EXISTS vehicles (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    chassis VARCHAR(100) NOT NULL,
    immat VARCHAR(100) NOT NULL
);

INSERT INTO vehicles (chassis, immat) VALUES ('CH1234XYZ', 'IM123456');
INSERT INTO vehicles (chassis, immat) VALUES ('CH5678XYZ', 'IM654321');
INSERT INTO vehicles (chassis, immat) VALUES ('CH9101XYZ', 'IM789012');
INSERT INTO vehicles (chassis, immat) VALUES ('CH1123XYZ', 'IM345678');
INSERT INTO vehicles (chassis, immat) VALUES ('CH1415XYZ', 'IM901234');
INSERT INTO vehicles (chassis, immat) VALUES ('CH1617XYZ', 'IM567890');
INSERT INTO vehicles (chassis, immat) VALUES ('CH1819XYZ', 'IM123098');
INSERT INTO vehicles (chassis, immat) VALUES ('CH2021XYZ', 'IM765432');
INSERT INTO vehicles (chassis, immat) VALUES ('CH2223XYZ', 'IM456789');
INSERT INTO vehicles (chassis, immat) VALUES ('CH2425XYZ', 'IM098765');