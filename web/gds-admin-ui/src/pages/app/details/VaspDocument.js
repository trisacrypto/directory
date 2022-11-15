import React from 'react';
import { Page, Text, View, Font, Document, StyleSheet, Link } from '@react-pdf/renderer';
import Html from 'react-pdf-html';
import { formatDate, formatDisplayedData, getConductsCustomerKYC, getMustComplyRegulations, getMustSafeguardPii, getSafeguardPii, hasAddressField, hasAddressFieldAndLine, hasAddressLine, intlFormatter, verifiedContactStatus } from 'utils';
import { BUSINESS_CATEGORY, NAME_IDENTIFIER_TYPE } from 'constants/basic-details';
import { AddressTypeHeaders, StatusLabel, VERIFIED_CONTACT_STATUS_LABEL } from 'constants/index';
import { financialTransfersPermitted, hasRequiredRegulatoryProgram } from 'constants/trixo';
import { isoCountries } from 'utils/country';
import { NATIONAL_IDENTIFIER_TYPE } from 'constants/national-identification';

Font.register({
    family: 'Source Sans Pro',
    fonts: [
        { src: 'https://fonts.gstatic.com/s/sourcesanspro/v14/6xK3dSBYKcSV-LCoeQqfX1RYOo3aPw.ttf' },
        { src: 'https://fonts.gstatic.com/s/sourcesanspro/v14/fpTVHK8qsXbIeTHTrnQH6EfrksRSinjQUrHtm_nW72g.ttf', fontStyle: 'italic', fontWeight: 600 },
        { src: 'https://fonts.gstatic.com/s/sourcesanspro/v14/6xKydSBYKcSV-LCoeQqfX1RYOo3i54rAkA.ttf', fontWeight: 600 },
    ],
    format: 'truetype'
});

// Create styles
const styles = StyleSheet.create({
    body: {
        paddingTop: 35,
        paddingBottom: 65,
        paddingHorizontal: 35,
        fontSize: 12,
        fontFamily: 'Source Sans Pro'
    },
    text: {
        fontSize: 12,
        marginBottom: 5,
        fontFamily: 'Source Sans Pro'
    },
    textBold: {
        fontWeight: 'bold'
    },
    header1: {
        marginTop: 30,
        marginBottom: 5,
        fontSize: 16,
        fontWeight: 'bold',
        fontFamily: 'Source Sans Pro'
    },
    header2: {
        marginTop: 15,
        marginBottom: 5,
        fontSize: 14,
        fontWeight: 'bold',
        fontFamily: 'Source Sans Pro'
    },
    dFlex: {
        display: 'flex',
        flexDirection: 'row',
        flexWrap: 'wrap'
    },
    bottomSpacing: {
        marginBottom: 5
    },
    textUnderline: {
        textDecoration: 'underline'
    }
});

const NameView = ({ name }) => {
    const names = React.useMemo(() => name && [].concat(...Object.values(name)), [name])

    return (
        <View>
            <Text style={styles.header2}>Name identifiers</Text>
            <View>
                {
                    names && names.map((value, index) => (
                        <View key={index}>
                            <Text>- {NAME_IDENTIFIER_TYPE[value?.legal_person_name_identifier_type]}: {formatDisplayedData(value?.legal_person_name)}</Text>
                        </View>
                    ))
                }
                <View style={[styles.dFlex, { alignItems: 'center' }]}></View>
            </View>
        </View>
    )
}

const AddressLines = ({ address }) => (
    <View>
        {
            address && address.address_line.map((addressLine, idx) => addressLine ? <Text key={idx}>{addressLine} {"\n"}</Text> : null)
        }
    </View>
)

const AddressField = ({ address }) => (
    <View>
        <Text>{address.sub_department ? <Text>{address?.sub_department}{"\n"}</Text> : null}</Text>
        <Text>{address.department ? address?.department : null}</Text>
        <Text>{address.building_number} {address?.street_name}</Text>
        <Text>{address.post_box ? <Text>P.O. Box: {address?.post_box}</Text> : null}</Text>
        <Text>{address.floor || address.room || address.building_name ? <Text>{address?.floor} {address?.room} {address?.building_name}</Text> : null}</Text>
        <Text>{address.district_name ? address?.district_name : null}</Text>
        <Text>{address.town_name || address.town_location_name || address.country_sub_division ? <Text>{address?.town_name} {address?.town_location_name} {address?.country_sub_division} {address?.post_code}</Text> : null}</Text>
        <Text>{address?.country}</Text>
    </View>
)

const Address = ({ address }) => {
    if (hasAddressFieldAndLine(address)) {
        return <Text>Invalid Address</Text>
    }

    if (hasAddressLine(address)) {
        return <AddressLines address={address} />
    }

    if (hasAddressField(address)) {
        return <AddressField address={address} />
    }

    return <Text style={{ fontStyle: 'italic' }}>Unparseable Address</Text>
}

const GeographicView = ({ data }) => {
    return (
        <View>
            <Text style={styles.header2}>Address(es)</Text>
            {
                data && data.map((address, index) => {
                    return (
                        <View key={index}>
                            <Text style={[styles.bottomSpacing, styles.textBold]}>{AddressTypeHeaders[address?.address_type]} Address:</Text>
                            <Address address={address} />
                        </View>
                    )
                })
            }
        </View>
    )
}

const TrisaDetailsView = ({ data }) => (
    <View style={{ lineHeight: 1.5 }}>
        <Text style={styles.header1}>TRISA Implementation</Text>
        <View>
            <Text style={styles.textBold}>ID: <Text style={{ fontWeight: 'normal' }}>{formatDisplayedData(data?.id)}</Text></Text>
            <Text style={styles.textBold}>Registered directory: <Text style={{ fontWeight: 'normal' }}>{formatDisplayedData(data?.registered_directory)}</Text></Text>
        </View>
        <View>
            <Text style={styles.textBold}>Common name: <Text style={{ fontWeight: 'normal' }}>{formatDisplayedData(data?.common_name)}</Text></Text>
            <Text style={styles.textBold}>Endpoint: <Text style={{ fontWeight: 'normal' }}>{formatDisplayedData(data?.trisa_endpoint)}</Text></Text>
        </View>
        <View >
            <Text style={styles.textBold}>First listed: <Text style={{ fontWeight: 'normal' }}>{formatDate(data?.first_listed)}</Text></Text>
            <Text style={styles.textBold}>Last updated: <Text style={{ fontWeight: 'normal' }}>{formatDate(data?.last_updated)}</Text></Text>
            <Text style={styles.textBold}>Verified on: <Text style={{ fontWeight: 'normal' }}>{formatDate(data?.verified_on)}</Text></Text>
        </View>
    </View>
)

const TrixoForm = ({ data }) => (
    <View>
        <Text style={styles.header1}>TRIXO Form</Text>
        <Text>Organization <Text style={styles.textBold}>{financialTransfersPermitted[data?.financial_transfers_permitted]} </Text>
            partially permitted to send and/or receive transfers of virtual assets in the jurisdiction(s) in which it operates.</Text>
        <View style={[styles.dFlex]}>
            <Text style={styles.textBold}>Primary National Jurisdiction: </Text>
            <Text>{isoCountries[data?.primary_national_jurisdiction]}</Text>
        </View>
        <View style={[styles.dFlex]}>
            <Text style={styles.textBold}>Name of Primary Regulator: </Text>
            <Text>{data?.primary_regulator}</Text>
        </View>
        <View>
            {
                Array.isArray(data?.other_jurisdictions) && data?.other_jurisdictions.map(jurisdiction => {
                    return (
                        <View>
                            <Text style={styles.header2}>Other Jurisdictions</Text>
                            <View style={[styles.dFlex]}>
                                <Text style={styles.textBold}>Country: </Text>
                                <Text>{isoCountries[jurisdiction?.country]}</Text>
                            </View>
                            <View>
                                <Text style={styles.textBold}>Regulator name: </Text>
                                <Text>{jurisdiction?.regulator_name}</Text>
                            </View>
                        </View>
                    )
                })
            }
        </View>


        <Text style={styles.header2}>CDD & Travel Rule Policies</Text>
        <Text style={{ marginTop: 3 }}>Organization <Text style={styles.textBold}>{hasRequiredRegulatoryProgram[data?.has_required_regulatory_program]}</Text> programme that sets minimum AML, CFT, KYC/CDD and Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes where it is licensed/approved/registered.</Text>
        <Text style={{ marginTop: 5 }}>Organization <Text style={styles.textBold}>{getConductsCustomerKYC(data?.conducts_customer_kyc)}</Text> conduct KYC/CDD before permitting its customers to send/receive virtual asset transfers.</Text>
        <Text style={styles.textBold}>Conducts KYC Threshold: <Text style={{ fontWeight: 'normal' }}>{intlFormatter({ currency: data?.kyc_threshold_currency }).format(data?.kyc_threshold)} {data?.kyc_threshold_currency}</Text> </Text>
        <Text style={{ marginTop: 5 }}>Organization <Text style={styles.textBold}>{getMustComplyRegulations(data?.must_comply_regulations)}</Text> comply with the application of the Travel Rule standards in the
            jurisdiction(s) where it is licensed/approved/registered.</Text>

        <View style={{ marginTop: 10 }}>
            <Text>
                Applicable Regulation(s):
            </Text>
            {
                data?.applicable_regulations.map(regulation => <Text key={regulation}>- {regulation}</Text>)
            }
        </View>
        <Text style={styles.textBold}>Minimum Compliance Threshold: <Text style={{ fontWeight: 'normal' }}>{intlFormatter({ currency: data?.compliance_threshold_currency }).format(data?.compliance_threshold)} {data?.compliance_threshold_currency}</Text></Text>

        <Text style={styles.header2}>Data Protection Policies</Text>
        <Text>Organization <Text style={styles.textBold}>{getMustSafeguardPii(data?.must_safeguard_pii)}</Text> required to safeguard PII by law.</Text>
        <Text style={{ marginTop: 5 }}>Organization <Text style={styles.textBold}>{getSafeguardPii(data?.safeguard_pii)}</Text> secure and protect PII, including PII received from other VASPs
            under the Travel Rule.</Text>
    </View>
)

const Contact = ({ data, type, verifiedContact }) => {
    const status = verifiedContactStatus({ data, type, verifiedContact })

    return (
        <View style={{ marginRight: 40, marginBottom: 15 }}>
            <Text style={styles.header2}>{type} contact:</Text>
            <View style={[{ lineHeight: 1.5 }]}>
                <Text>{formatDisplayedData(data?.name)}</Text>
                <Text>{formatDisplayedData(data?.phone)}</Text>
                <View style={[styles.dFlex]}>
                    <Text>{formatDisplayedData(data?.email)} </Text>
                    <Text style={{ fontStyle: 'italic' }}>{VERIFIED_CONTACT_STATUS_LABEL[status]}</Text>
                </View>
                <Text>{data?.person ? 'Has IVMS101 Record' : 'No IVMS101 Data'}</Text>
            </View>
        </View>
    )
}

const ContactView = ({ data, verifiedContact }) => (
    <View>
        <Text style={styles.header1}>Contacts</Text>
        <View style={[styles.dFlex]}>
            {data?.technical ? <Contact verifiedContact={verifiedContact} data={data?.technical} type="Technical"></Contact> : null}
            {data?.administrative ? <Contact verifiedContact={verifiedContact} data={data?.administrative} type="Administrative" /> : null}
            {data?.legal ? <Contact verifiedContact={verifiedContact} data={data?.legal} type="Legal" /> : null}
            {data?.billing ? <Contact verifiedContact={verifiedContact} data={data?.billing} type="Billing" /> : null}
        </View>
    </View>
)

const NationalIdentificationView = ({ data }) => (
    <View>
        <Text style={styles.header2}>National Identification:</Text>
        <View style={[styles.dFlex]}>
            <Text style={styles.textBold}>Type: </Text>
            <Text>LEI </Text>
            <Text style={styles.textBold}>LEIX: </Text>
            <Text>{formatDisplayedData(data?.national_identifier)}</Text>
        </View>
        <Text>Issued by: {`(${formatDisplayedData(data?.country_of_issue)}) by authority ${formatDisplayedData(data?.registration_authority)}`}</Text>
        <View style={[styles.dFlex, { alignItems: 'center' }]}>
            <Text>National identification type:</Text>
            <Text>{formatDisplayedData(NATIONAL_IDENTIFIER_TYPE[data?.national_identifier_type])}</Text>
        </View>
    </View>
)

const CertificateDetailsView = ({ data }) => (
    <View>
        <Text style={styles.header2}>Certificate Details</Text>
        <View>
            <Text style={styles.textBold}>Expires: <Text style={{ fontWeight: 'normal' }}>{new Date(data?.not_after).toUTCString()}</Text></Text>
            <Text style={styles.textBold}>Serial Number: <Text style={{ fontWeight: 'normal' }}>{formatDisplayedData(data?.serial_number)}</Text></Text>
            <Text style={styles.textBold}>Issuer: <Text style={{ fontWeight: 'normal' }}></Text>{data?.issuer?.common_name}</Text>
        </View>
    </View>
)

const ReviewNotes = ({ notes }) => {
    return notes && <View>
        <Text style={styles.header1}>Reviewer note(s)</Text>
        {
            notes.map((note, idx) => <ReviewNote note={note} key={idx} />)
        }
    </View>
}

const ReviewNote = ({ note }) => {
    return (
        <View style={{ marginBottom: 5 }}>
            {note?.id ?
                <View style={[styles.dFlex]}>
                    <Text>
                        <Html style={{ fontSize: 12 }}>{note?.text}</Html>
                        <Text style={{ fontStyle: 'italic' }}>{note?.editor ? '(edited)' : ''}</Text></Text>
                </View> : null
            }
            {note?.id ? <View>
                {
                    note?.editor ?
                        <Text>
                            by <Link>{note?.editor}</Link> on {new Date(note?.modified).toUTCString()}
                        </Text> : <Text>
                            by <Link>{note?.author}</Link> on {new Date(note?.created).toUTCString()}
                        </Text>
                }

            </View> : null}
        </View>
    )
}

// Create Document Component
const VaspDocument = ({ vasp, notes }) => {

    return (
        <Document wrap={false}>
            <Page size="A4" style={styles.body}>
                <View>
                    <View>
                        <View>
                            <Text style={styles.header1}>{vasp?.name}</Text>
                            <View style={styles.dFlex}>
                                {vasp ? <Text>{StatusLabel[vasp?.vasp?.verification_status]}, </Text> : null}
                                {vasp?.traveler ? <Text>Traveler</Text> : <Text>Not Traveler</Text>}
                            </View>
                            <Link>{vasp?.vasp?.website}</Link>

                            <View style={[styles.dFlex, { marginTop: 10 }]}>
                                <Text style={styles.textBold}>Customer Number: </Text>
                                <Text>{formatDisplayedData(vasp?.vasp?.entity?.customer_number)}</Text>
                            </View>
                            <View style={[styles.dFlex]}>
                                <Text style={styles.textBold}>Country of Registration: </Text>
                                <Text>{formatDisplayedData(vasp?.vasp?.entity?.country_of_registration)}</Text>
                            </View>
                        </View>

                        <Text style={styles.header1}>Business details</Text>
                        <View style={styles.dFlex}>
                            <Text style={styles.textBold}>Established on: </Text>
                            <Text>{formatDisplayedData(vasp?.vasp?.established_on)}</Text>
                        </View>
                        <View style={[styles.dFlex]}>
                            <Text style={styles.textBold}>Business Category: </Text>
                            <Text>{BUSINESS_CATEGORY[vasp?.vasp?.business_category]}</Text>
                        </View>
                        <View>
                            <Text style={styles.header2}>VASP categorie(s):</Text>
                            <View>
                                {
                                    vasp?.vasp?.vasp_categories && vasp?.vasp?.vasp_categories.map((category, index) => <Text key={index} >  - {category}</Text>)
                                }
                            </View>
                        </View>
                        <NameView name={vasp?.vasp?.entity?.name} />
                        <GeographicView data={vasp?.vasp?.entity?.geographic_addresses} />
                        <NationalIdentificationView data={vasp?.vasp?.entity?.national_identification} />
                        <ContactView data={vasp?.vasp?.contacts} verifiedContact={vasp?.verified_contacts} />
                    </View>
                    <TrixoForm data={vasp?.vasp?.trixo} />
                    <TrisaDetailsView data={vasp?.vasp} />
                    <CertificateDetailsView data={vasp?.vasp?.identity_certificate} />
                    <View>
                        <ReviewNotes notes={notes} />
                    </View>
                </View>
            </Page>
        </Document>
    )
};


export default VaspDocument;