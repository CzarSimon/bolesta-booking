import { Response } from "express";
import { StatusOK, StatusNotFound } from "../constants";

export function sendOK(res: Response) {
    res.status(StatusOK).json({status: "OK"});
}

export function sendNotFound(res: Response) {
    res.status(StatusNotFound).json({status: "NOT FOUND"});
}