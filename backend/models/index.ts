export interface Cabin {
    id: string;
    name: string;
    createdAt: Date;
    updatedAt: Date;
}

export interface User {
    id: string;
    name: string;
    email: string;
    createdAt: Date;
    updatedAt: Date;
}

export interface Booking {
    id: string;
    startDate: Date;
    endDate: Date;
    createdAt: Date;
    updatedAt: Date;
    cabin: Cabin;
    user: User;
}

export interface BookingRequest {
    cabinId: string;
    startDate: Date;
    endDate: Date;
    userId: string;
    password: string;
};

export interface BookingFilter {
    cabinId?: string;
    userId?: string;
}