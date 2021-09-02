// @flow
import React from 'react';
import { Link } from 'react-router-dom';
import { Row, Col, Card, Button } from 'react-bootstrap';
import classNames from 'classnames';

import PageTitle from '../../../components/PageTitle';
import Table from '../../../components/Table';
import { useDispatch, useSelector } from 'react-redux';

import { fetchCertificates, fetchVasps } from '../../../redux/dashboard/actions';
import { Status } from '../../../constants';

const NameColumn = ({ row }) => {
    return (
        <React.Fragment>
            <p className="m-0 d-inline-block align-middle font-16">
                <Link to="/" className="text-body">
                    {row.original.common_name}
                </Link>
            </p>
        </React.Fragment>
    );
};

const StatusColumn = ({ row }) => {
    return (
        <React.Fragment>
            <span
                className={classNames('badge', {
                    'bg-success': row.original.verification_status === "VERIFIED",
                    'bg-warning': row.original.verification_status === "SUBMITTED" || row.original.verification_status === "PENDING_REVIEW",
                })}>
                {Status[row.original.verification_status]}
            </span>
        </React.Fragment>
    );
};

const ActionColumn = ({ row }) => {
    const id = row?.original?.id || "";

    return (
        <React.Fragment>
            <Link to={`/vasps-summary/${id}/details`} className="action-icon text-center">
                <i className="mdi mdi-eye"></i>
            </Link>
        </React.Fragment>
    );
};

const columns = [
    {
        Header: 'Name',
        accessor: 'name',
        sort: true,
        Cell: NameColumn,
    },
    {
        Header: 'Status',
        accessor: 'verification_status',
        sort: true,
        Cell: StatusColumn
    },
    {
        Header: 'Established On',
        accessor: 'established_on',
        sort: true,
    },
    {
        Header: 'Action',
        accessor: 'action',
        sort: false,
        classes: 'table-action',
        Cell: ActionColumn,
    },
];


const VaspsList = (): React$Element<React$FragmentType> => {
    const dispatch = useDispatch()
    const vasps = useSelector((state) => state.Vasps.data)

    const sizePerPageList = [
        {
            text: '5',
            value: 5,
        },
        {
            text: '10',
            value: 10,
        },
        {
            text: '20',
            value: 20,
        }
    ];

    React.useEffect(() => {
        dispatch(fetchVasps())
        dispatch(fetchCertificates())
    }, [dispatch])

    return (
        <React.Fragment>
            <PageTitle
                breadCrumbItems={[]}
                title={'VASPs list'}
            />

            <Row>
                <Col>
                    <Card>
                        <Card.Body>
                            <Row>
                                <Col sm={12}>
                                    <div className="text-sm-end">
                                        <Button className="btn btn-success mb-2 me-1">
                                            <i className="mdi mdi-cog-outline"></i>
                                        </Button>

                                        <Button className="btn btn-light mb-2 me-1">Import</Button>

                                        <Button className="btn btn-light mb-2">Export</Button>
                                    </div>
                                </Col>
                            </Row>
                            {
                                vasps && vasps.length &&
                                <Table
                                    columns={columns}
                                    data={vasps}
                                    pageSize={5}
                                    sizePerPageList={sizePerPageList}
                                    isSortable={true}
                                    pagination={true}
                                    isSelectable={true}
                                    isSearchable={true}
                                    theadClass="table-light"
                                    searchBoxClass="mt-2 mb-3"
                                />
                            }
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
        </React.Fragment>
    );
};

export default VaspsList;
