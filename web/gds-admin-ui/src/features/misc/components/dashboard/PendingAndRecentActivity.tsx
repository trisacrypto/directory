import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import React from 'react';
import { Card, Dropdown, Table } from 'react-bootstrap';
import { useHistory } from 'react-router-dom';
import SimpleBar from 'simplebar-react';

import CiphertraceFavicon from '@/assets/images/ciphertrace.ico';
import TrisaFavicon from '@/assets/images/trisa_favicon.png';
import NoData from '@/components/NoData';
import { StatusLabel } from '@/constants/index';
import { actionType, useModal } from '@/contexts/modal';
import { useGetVasps } from '@/features/vasps';

dayjs.extend(relativeTime);

const PENDING_REVIEW_QUERY_PARAMS = 'status=pending_review';

const PendingReviewsTable = ({ data }: any) => {
    const [vasp, setVasp] = React.useState({ name: '', id: '' });
    const { dispatch } = useModal();
    const history = useHistory();

    const handleResendEmailClick = (name: any) => {
        setVasp(name);
        dispatch({ type: actionType.SEND_EMAIL_MODAL, payload: { vasp } });
    };

    return (
        <SimpleBar style={{ maxHeight: 350 }} className="task">
            <Table responsive className="table table-centered table-nowrap table-hover mb-0 z-index-2">
                <tbody>
                    {data?.map((vasp: any) => (
                        <tr key={vasp.id}>
                            <td
                                onClick={() => history.push(`/vasps/${vasp?.id}`)}
                                className="d-flex gap-2 align-items-center"
                                role="button">
                                <div>
                                    {vasp?.traveler ? (
                                        <img src={CiphertraceFavicon} width="30" alt="CipherTrace" />
                                    ) : (
                                        <img src={TrisaFavicon} width="30" className="img-fluid" alt="TRISA" />
                                    )}
                                </div>
                                <div>
                                    <h5 className="font-14 my-1 gap-2 d-flex">{vasp?.name}</h5>
                                    <span className="text-muted font-13">{dayjs(vasp?.last_updated).fromNow()}</span>
                                </div>
                            </td>
                            <td>
                                <span className="text-muted font-13">Status</span> <br />
                                <span className="badge badge-warning-lighten">
                                    {(StatusLabel as any)[vasp?.verification_status]}
                                </span>
                            </td>
                            <td>
                                <span className="text-muted font-13">Emails</span>
                                <div className="font-14 mt-1 fw-normal">
                                    {vasp?.verified_contacts &&
                                        Object.keys(vasp?.verified_contacts).map((contact) => (
                                            <span
                                                key={contact}
                                                className={`badge ${
                                                    vasp?.verified_contacts[contact]
                                                        ? 'badge-success-lighten'
                                                        : 'badge-danger-lighten'
                                                }`}
                                                style={{ marginRight: 4 }}>
                                                {contact}
                                            </span>
                                        ))}
                                </div>
                            </td>
                            <td className="table-action text-center" style={{ width: '90px' }}>
                                <Dropdown align="end">
                                    <Dropdown.Toggle variant="link" className="arrow-none card-drop p-0 shadow-none">
                                        <i className="mdi mdi-dots-horizontal" />
                                    </Dropdown.Toggle>
                                    <Dropdown.Menu>
                                        <Dropdown.Item onClick={() => history.push(`/vasps/${vasp?.id}`)}>
                                            {' '}
                                            <span className="mdi mdi-eye-outline" /> View
                                        </Dropdown.Item>
                                        <Dropdown.Item onClick={() => handleResendEmailClick(vasp)}>
                                            {' '}
                                            <span className="mdi mdi-email-edit-outline" /> Resend email
                                        </Dropdown.Item>
                                    </Dropdown.Menu>
                                </Dropdown>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </Table>
        </SimpleBar>
    );
};

const PendingAndRecentActivity = () => {
    const { data } = useGetVasps({
        queryParams: PENDING_REVIEW_QUERY_PARAMS,
    });

    return (
        <Card className="w-100">
            <Card.Body>
                <h4 className="header-title mb-3">Pending and Recent Activity</h4>
                {data?.vasps?.length ? (
                    <PendingReviewsTable data={data?.vasps} />
                ) : (
                    <NoData title="No available pending registrations" />
                )}
            </Card.Body>
        </Card>
    );
};

export default PendingAndRecentActivity;
