// @flow
import React from 'react';
import { Link } from 'react-router-dom';
import { Row, Col, Card, Button } from 'react-bootstrap';
import classNames from 'classnames';
import dayjs from 'dayjs';

import relativeTime from 'dayjs/plugin/relativeTime'


import PageTitle from '../../../components/PageTitle';
import Table from '../../../components/Table';
import { useDispatch, useSelector } from 'react-redux';

import { fetchCertificates, fetchVasps } from '../../../redux/dashboard/actions';
import { Status, StatusLabel } from '../../../constants';
dayjs.extend(relativeTime)


const NameColumn = ({ row }) => {
    const id = row?.original?.id || "";

    return (
        <React.Fragment>
            <p className="m-0 d-inline-block align-middle font-16">
                <Link to={`/vasps/${id}/details`} className="text-body">
                    {row.original.name}
                    <span className="text-muted font-italic d-block">
                        {row.original.common_name}
                    </span>
                </Link>
            </p>
        </React.Fragment>
    );
};

const LastUpdatedColumn = ({ row }) => {

    return <React.Fragment>
        <p className="m-0 d-inline-block align-middle font-16">
            <span>
                {dayjs(row?.original?.last_updated).fromNow()}
            </span>
        </p>
    </React.Fragment>
}

const StatusColumn = ({ row }) => {
    return (
        <React.Fragment>
            <span
                className={classNames('badge', {
                    'bg-success': row.original.verification_status === Status.VERIFIED,
                    'bg-warning': row.original.verification_status === Status.SUBMITTED || row.original.verification_status === Status.PENDING_REVIEW,
                })}>
                {StatusLabel[row.original.verification_status]}
            </span>
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
        Header: 'Last updated',
        accessor: 'last_updated',
        Cell: LastUpdatedColumn,
        sort: true,
    }
];


const VaspsList = (): React$Element<React$FragmentType> => {
    const dispatch = useDispatch()
    const data = useSelector((state) => state.Vasps.data)

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
        },
        {
            text: '100',
            value: 100
        }
    ];

    React.useEffect(() => {
        dispatch(fetchVasps())
        dispatch(fetchCertificates())
    }, [dispatch])

    return (
        <React.Fragment>
            <PageTitle
                breadCrumbItems={[
                    { label: 'List', path: '/vasps', active: true }
                ]}
                title={'All Registered VASPs'}
            />

            <Row>
                <Col>
                    <Card>
                        <Card.Body>
                            <Row>
                                <Col sm={12}>
                                    <div className="text-sm-end">
                                        <Button className="btn btn-light mb-2">Export</Button>
                                    </div>
                                </Col>
                            </Row>
                            {
                                data && data.vasps && data.vasps.length &&
                                <Table
                                    columns={columns}
                                    data={data?.vasps}
                                    pageSize={data.page_size || 100}
                                    sizePerPageList={sizePerPageList}
                                    isSortable={true}
                                    pagination={true}
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
