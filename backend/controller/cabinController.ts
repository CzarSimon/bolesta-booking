import { Request, Response } from "express";
import { METHOD, StatusOK } from "../constants";
import { sendNotFound } from "../httputil/response";
import { Server } from "../httputil/server";
import { CabinService } from "../service";

export class CabinController {
  private service: CabinService;

  constructor(service: CabinService) {
    this.service = service;
  }

  public attach(server: Server): void {
    server.register(METHOD.GET, "/v1/cabins", this.listCabins);
    server.register(METHOD.GET, "/v1/cabins/:id", this.getCabin);
  }

  private getCabin = async (req: Request, res: Response) => {
    const cabinId = req.params.id;
    const cabin = await this.service.getCabin(cabinId);
    if (!cabin) {
      console.log(`[ERROR]: Cabin with id=${cabinId} was not found`);
      sendNotFound(res);
      return;
    }

    res.status(StatusOK).json(cabin);
  };

  private listCabins = async (req: Request, res: Response) => {
    const cabins = await this.service.listCabins();
    res.status(StatusOK).json(cabins);
  };
}
