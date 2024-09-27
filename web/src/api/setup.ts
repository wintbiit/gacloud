import {axioser} from "./index.ts";

export interface MySQLOptionsProps {
    type: string,
    host: string,
    port: number,
    username: string,
    password: string,
    database: string,
}

export interface PostgreSQLOptionsProps {
    type: string
    host: string,
    port: number,
    username: string,
    password: string,
    database: string,
}

export interface SQLiteOptionsProps {
    type: string,
    path: string,
}

interface SetupStatus {
    currentStep: number,
}

async function getSetupStatus() : Promise<number> {
    return axioser.get<SetupStatus>("/setup").then(res => res.data.currentStep);
}

const setupDatabase = async (values: MySQLOptionsProps | PostgreSQLOptionsProps | SQLiteOptionsProps) => {
    await axioser.post("/setup/database", {
        type: values.type,
        params: values,
    });
}

interface TestDatabaseResponse {
    success: boolean,
}

async function testDatabase(values: MySQLOptionsProps | PostgreSQLOptionsProps | SQLiteOptionsProps) : Promise<boolean> {
    return axioser.post<TestDatabaseResponse>("/setup/database/test", {
        type: values.type,
        params: values,
    }).then(res => res.data.success);
}

export {
    getSetupStatus,
    setupDatabase,
    testDatabase
}