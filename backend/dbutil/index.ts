import sqlite3, { RunResult } from "sqlite3";
import { SCHEMA } from "./schema";

export class Database extends sqlite3.Database {
  constructor(name: string, runMigrations: boolean = true) {
    super(name);

    if (runMigrations) {
      migrate(this);
    }
  }

  public async runQuery<T = void>(sql: string, params?: T): Promise<RunResult> {
    const stmt = this.prepare(sql);
    return new Promise((resolve, reject) => {
      stmt.run(params, (err: Error, res: RunResult) => {
        if (err) {
          reject(err);
        }
        resolve(res);
      });
    });
  }

  public async queryOne<T>(sql: string, params: T): Promise<any> {
    const stmt = this.prepare(sql);

    return new Promise((resolve, reject) => {
      stmt.get(params, (err, row) => {
        if (err) {
          reject(err);
        }
        resolve(row);
      });
    });
  }

  public async queryAll<T = void>(sql: string, params?: T): Promise<any[]> {
    const stmt = this.prepare(sql);

    return new Promise((resolve, reject) => {
      stmt.all(params, (err, rows) => {
        if (err) {
          console.log(err);
          reject(err);
        }
        resolve(rows);
      });
    });
  }
}

async function migrate(db: Database) {
  await db.runQuery("PRAGMA foreign_keys = ON;");
  for (const sql of SCHEMA) {
    await db.runQuery(sql);
  }
}
