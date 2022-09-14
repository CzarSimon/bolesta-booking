import { Cabin } from "../models";
import { now } from "../timeutil";
import { Optional } from "../types";

export interface CabinRepository {
    find(id: string): Optional<Cabin>;
    findAll(): Cabin[]
}

export function newCabinRepository() {
    return new DummyCabinRepository([
        {
            id: "a4b4f496-767e-423e-9816-83b71e1cfa89",
            name: "Bölestastugan",
            createdAt: now(),
            updatedAt: now(),
        }
    ])
}

class DummyCabinRepository implements CabinRepository {
    cabins: Record<string, Cabin>

    constructor(cabins: Cabin[]) {
        this.cabins = {};
        cabins.forEach(cabin => {
            this.cabins[cabin.id] = cabin;
        });
    }

    public find(id: string): Optional<Cabin> {
        return this.cabins[id]
    }

    public findAll(): Cabin[] {
        return Object.values(this.cabins);
    }
}