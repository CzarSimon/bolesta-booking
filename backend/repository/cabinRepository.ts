import { Database } from "../dbutil";
import { Cabin } from "../models";
import { Optional } from "../types";

export interface CabinRepository {
  find(id: string): Promise<Optional<Cabin>>;
  findAll(): Promise<Cabin[]>;
}

export function newCabinRepository(db: Database) {
  return new SQLiteCabinRepository(db);
}

class SQLiteCabinRepository implements CabinRepository {
  db: Database;

  constructor(db: Database) {
    this.db = db;
  }

  public async find(id: string): Promise<Optional<Cabin>> {
    const res = await this.db.queryOne<string>(
      "SELECT id, name, created_at, updated_at FROM cabin WHERE id = ?",
      id
    );

    if (!res) {
      return undefined;
    }

    return parseCabin(res);
  }

  public async findAll(): Promise<Cabin[]> {
    const res = await this.db.queryAll(
      "SELECT id, name, created_at, updated_at FROM cabin"
    );

    return res.map((row) => parseCabin(row));
  }
}

function parseCabin(row: any): Cabin {
  return {
    id: row.id,
    name: row.name,
    createdAt: new Date(row.created_at),
    updatedAt: new Date(row.updated_at),
  };
}
