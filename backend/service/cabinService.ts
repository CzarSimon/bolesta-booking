import { Cabin } from "../models";
import { Optional } from "../types";
import { CabinRepository } from "../repository";

export class CabinService {
  cabinRepository: CabinRepository;

  constructor(cabinRepository: CabinRepository) {
    this.cabinRepository = cabinRepository;
  }

  public async getCabin(id: string): Promise<Optional<Cabin>> {
    return this.cabinRepository.find(id);
  }

  public listCabins(): Promise<Cabin[]> {
    return this.cabinRepository.findAll();
  }
}
