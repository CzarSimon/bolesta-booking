import { Request, Response } from "express";
import { METHOD, StatusOK } from "../constants";
import { Server } from "../httputil/server";
import { BookingFilter, BookingRequest } from "../models";
import { BookingService } from "../service";

export class BookingController {
  private service: BookingService;

  constructor(service: BookingService) {
    this.service = service;
  }

  public attach(server: Server): void {
    server.register(METHOD.GET, "/v1/bookings", this.listBookins);
    server.register(METHOD.POST, "/v1/bookings", this.createBooking);
  }

  private listBookins = (
    req: Request<{}, {}, {}, BookingFilter>,
    res: Response
  ): void => {
    const { cabinId, userId } = req.query;
    const bookings = this.service.listBookings({ cabinId, userId });
    res.status(StatusOK).json(bookings);
  };

  private async createBooking(req: Request, res: Response) {
    const bookingRequest = parseBookingRequest(req);
    const booking = await this.service.createBooking(bookingRequest);
    res.status(StatusOK).json(booking);
  }
}

function parseBookingRequest(req: Request): BookingRequest {
  const { cabinId, userId, startDate, endDate, password } = req.body;
  return {
    cabinId,
    startDate: new Date(startDate),
    endDate: new Date(endDate),
    userId,
    password,
  };
}
