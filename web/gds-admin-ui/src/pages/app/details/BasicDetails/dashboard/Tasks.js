/* eslint-disable no-unused-vars */
import React from 'react';
import { Card, Dropdown, Table } from 'react-bootstrap';
import { useHistory } from 'react-router-dom'
import dayjs from 'dayjs';
import { StatusLabel, Status as VerificationStatus } from 'constants/index';
import TrisaFavicon from 'assets/images/trisa_favicon.png'
import CiphertraceFavicon from 'assets/images/ciphertrace.ico'
import SimpleBar from 'simplebar-react';
import PropTypes from 'prop-types';
import relativeTime from 'dayjs/plugin/relativeTime'
import { actionType, useModal } from 'contexts/modal';
import OvalLoader from 'components/OvalLoader';
import NoData from 'components/NoData';
import filter from 'lodash/filter';
import isDate from 'lodash/isDate';
dayjs.extend(relativeTime)

/**
 * Verifiy that the passed date is less than 30days
 * @param {date} date 
 * @returns 
 */
function isRecent(date = '') {

    if (dayjs(date).isValid) {
        const tirthyDaysInMs = 30 * 24 * 60 * 60 * 1000;
        const now = dayjs();
        const dateDiff = now.diff(date)

        return Math.abs(dateDiff) <= tirthyDaysInMs
    }

    return false
}

export const getRecentVasps = (vasps) => vasps?.filter(vasp => {
    const now = dayjs()
    const certificateExpirationDate = vasp?.certificate_expiration && dayjs(vasp?.certificate_expiration)
    const certificateIssuedDate = vasp?.certificate_issued && dayjs(vasp?.certificate_issued)
    const verificationStatus = vasp?.verification_status

    if (verificationStatus === VerificationStatus.PENDING_REVIEW) {
        return vasp
    }
    if (certificateExpirationDate && certificateIssuedDate) {

        if (dayjs(certificateExpirationDate).isValid() || dayjs(certificateIssuedDate).isValid()) {

            if (certificateExpirationDate.isAfter(now) && isRecent(certificateExpirationDate)) {
                return vasp
            }

            if (certificateIssuedDate.isBefore(now) && isRecent(certificateIssuedDate)) {
                return vasp
            }
        }
    }
})

const PendingReviewsTable = ({ data }) => {
    const [vasp, setVasp] = React.useState({ name: '', id: '' });
    const { dispatch } = useModal()
    const history = useHistory()

    const handleResendEmailClick = (name) => {
        setVasp(name)
        dispatch({ type: actionType.SEND_EMAIL_MODAL, payload: { vasp } })
    }

    return (
        <SimpleBar style={{ maxHeight: 350 }} className="task">
            <Table responsive className="table table-centered table-nowrap table-hover mb-0 z-index-2">
                <tbody>
                    {
                        data?.map(vasp => (
                            <tr key={vasp.id}>
                                <td onClick={() => history.push(`/vasps/${vasp?.id}`)} className="d-flex gap-2 align-items-center" role="button">
                                    <div>
                                        {
                                            vasp?.traveler ? <img src={CiphertraceFavicon} width="30" alt="CipherTrace" /> : <img src={TrisaFavicon} width="30" className="img-fluid" alt="TRISA" />
                                        }
                                    </div>
                                    <div>
                                        <h5 className="font-14 my-1 gap-2 d-flex">
                                            {vasp?.name}
                                        </h5>
                                        <span className="text-muted font-13">{dayjs(vasp?.last_updated).fromNow()}</span>
                                    </div>
                                </td>
                                <td>
                                    <span className="text-muted font-13">Status</span> <br />
                                    <span className="badge badge-warning-lighten">{StatusLabel[vasp?.verification_status]}</span>
                                </td>
                                <td>
                                    <span className="text-muted font-13">Emails</span>
                                    <div className="font-14 mt-1 fw-normal">
                                        {
                                            vasp?.verified_contacts && Object.keys(vasp?.verified_contacts).map(contact => (
                                                <span key={contact} className={`badge ${vasp?.verified_contacts[contact] ? "badge-success-lighten" : "badge-danger-lighten"}`} style={{ marginRight: 4 }}>{contact}</span>
                                            )
                                            )
                                        }
                                    </div>
                                </td>
                                <td className="table-action text-center" style={{ width: '90px' }}>
                                    <Dropdown align="end">
                                        <Dropdown.Toggle variant="link" className="arrow-none card-drop p-0 shadow-none">
                                            <i className="mdi mdi-dots-horizontal"></i>
                                        </Dropdown.Toggle>
                                        <Dropdown.Menu>
                                            <Dropdown.Item onClick={() => history.push(`/vasps/${vasp?.id}`)}> <span className="mdi mdi-eye-outline"></span> View</Dropdown.Item>
                                            <Dropdown.Item onClick={() => handleResendEmailClick(vasp)}> <span className="mdi mdi-email-edit-outline"></span> Resend email</Dropdown.Item>
                                        </Dropdown.Menu>
                                    </Dropdown>
                                </td>
                            </tr>)
                        )
                    }
                </tbody>
            </Table>
        </SimpleBar>
    )
}

const Tasks = ({ data, isLoading }) => {
    const recentVasps = getRecentVasps(data?.vasps || [])

    if (isLoading || !data) {
        return (
            <Card className='w-100'>
                <Card.Body>
                    <h4 className="header-title mb-3">Pending and Recent Activity</h4>
                    <OvalLoader />
                </Card.Body>
            </Card>
        )
    }

    return (
        <Card className='w-100'>
            <Card.Body>
                <h4 className="header-title mb-3">Pending and Recent Activity</h4>
                {
                    (isLoading || !recentVasps) && <OvalLoader />
                }
                {
                    recentVasps?.length ? <PendingReviewsTable data={recentVasps} /> : <NoData title='No available pending registrations' />
                }
            </Card.Body>
        </Card>
    );
};

Tasks.propTypes = {
    data: PropTypes.object
}

export default Tasks;
