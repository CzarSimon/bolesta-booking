import { Response } from "express";
import { StatusOK } from "../constants";

export function sendOK(res: Response) {
    res.status(StatusOK).json({status: "OK"});
}