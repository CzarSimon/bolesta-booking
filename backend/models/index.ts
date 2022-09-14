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
    startDate: string;
    endDate: string;
    createdAt: Date;
    updatedAt: Date;
    cabin: Cabin;
    user: User;
}