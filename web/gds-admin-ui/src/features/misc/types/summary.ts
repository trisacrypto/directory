export type Summary = {
    vasps_count: number;
        pending_registrations: number;
        contacts_count: number;
        verified_contacts: number;
        certificates_issued: number;
        statuses: {
            REJECTED: number;
            SUBMITTED: number;
            VERIFIED: number;
        };
        certreqs: {
            COMPLETED: number;
            DOWNLOADED: number;
            INITIALIZED: number;
            READY_TO_SUBMIT: number;
        };
}