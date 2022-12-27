export interface ITrixo {
    primary_national_jurisdiction: string;
    primary_regulator: string;
    other_jurisdictions: string[];
    financial_transfers_permitted: string;
    has_required_regulatory_program: string;
    conducts_customer_kyc: boolean;
    kyc_threshold: number;
    kyc_threshold_currency: string;
    must_comply_travel_rule: boolean;
    applicable_regulations: string[];
    compliance_threshold: number;
    compliance_threshold_currency: string;
    must_safeguard_pii: boolean;
    safeguards_pii: boolean;
}

