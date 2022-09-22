import React from 'react';
import { Link, useHistory, useLocation } from 'react-router-dom';
import { Row, Col, Card, Button } from 'react-bootstrap';
import classNames from 'classnames';
import dayjs from 'dayjs';
import Select from 'react-select'
import qs from 'query-string'
import relativeTime from 'dayjs/plugin/relativeTime'


import PageTitle from 'components/PageTitle';
import Table from 'components/Table';
import { useDispatch, useSelector } from 'react-redux';

import { fetchVasps } from '../../../redux/dashboard/actions';
import { StatusLabel } from '../../../constants';
import { exportToCsv, getStatusClassName } from '../../../utils';
import useSafeDispatch from 'hooks/useSafeDispatch';
import { getAllVasps, getVaspsLoadingState } from 'redux/selectors/vasps';
import OvalLoader from 'components/OvalLoader';
dayjs.extend(relativeTime)


export const NameColumn = ({ row }) => {
    const id = row?.original?.id || "";

    return (
        <React.Fragment>
            <p className="m-0 d-inline-block align-middle font-16">
                <Link to={`/vasps/${id}`} className="text-body">
                    <span data-testid="name">
                        {row.original.name || 'N/A'}
                    </span>
                    <span className="text-muted font-italic d-block" data-testid="commonName">
                        {row.original.common_name || 'N/A'}
                    </span>
                </Link>
            </p>
        </React.Fragment>
    );
};

export const CertificateExpirationColumn = ({ row }) => {

    return (
        <React.Fragment>
            <p className="m-0 d-inline-block align-middle font-16" data-testid="certificate_expiration">
                {row?.original?.certificate_expiration ? dayjs(row?.original?.certificate_expiration).format("MMM DD, YYYY h:mm:ss a") : 'N/A'}
            </p>
        </React.Fragment>
    );
};


export const LastUpdatedColumn = ({ row }) => {

    return <React.Fragment>
        <p className="m-0 d-inline-block align-middle font-16">
            <span data-testid="last_updated">
                {row?.original?.last_updated ? dayjs(row?.original?.last_updated).fromNow() : 'N/A'}
            </span>
        </p>
    </React.Fragment>
}

export const StatusColumn = ({ row }) => {
    return (
        <React.Fragment>
            <span
                data-testid="verification_status"
                className={classNames('badge', getStatusClassName(row.original.verification_status))}>
                {StatusLabel[row.original.verification_status] || 'N/A'}
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
    },
    {
        Header: 'Certificate Expiration',
        accessor: 'certificate_expiration',
        sort: true,
        Cell: CertificateExpirationColumn
    }
];

const options = () => Object.entries(StatusLabel).map(([k, v]) => ({ value: k, label: v }))

const getOption = (key) => {
    const _options = options()
    let opt = []

    for (let option of _options) {
        if (typeof key === 'string' && option.value.toLowerCase() === key) {
            opt = [option]
        }

        if (Array.isArray(key)) {
            for (let k of key) {
                if (k === option.value.toLowerCase()) {
                    opt.push(option)
                }
            }

        }
    }
    return opt
}

const customStyles = {
    control: (styles) => ({ ...styles, paddingLeft: '9px !important' })
};


const VaspsList = () => {
    const location = useLocation()
    const parsedQuery = qs.parse(location.search)
    const query = qs.stringify(parsedQuery);

    const [queryParams, setQueryParams] = React.useState(query)
    const dispatch = useDispatch()
    const safeDispatch = useSafeDispatch(dispatch)
    const data = useSelector(getAllVasps)
    const isLoading = useSelector(getVaspsLoadingState)
    const history = useHistory()
    const [selectedRows, setSelectedData] = React.useState([]);

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
        safeDispatch(fetchVasps({ queryParams }))
    }, [safeDispatch, queryParams])

    const getQueryString = (arr) => {
        const queries = arr && Array.isArray(arr) ? arr.map(v => v.value.toLowerCase()) : ''
        return qs.stringify({ status: queries })
    }

    const handleSelectChange = (value) => {
        const params = getQueryString(value)
        setQueryParams(params)
        history.push({
            pathname: '/vasps',
            search: params
        })
    }

    const handleCsvExportClick = (rows) => selectedRows.length ? exportToCsv(selectedRows) : exportToCsv(rows)

    const onSelectedRows = rows => {
        const mappedRows = rows.map(r => r.original);
        setSelectedData([...mappedRows]);
    };

    const getExportButtonLabel = () => {
        if (selectedRows.length === 1) {
            return 'Export 1 row'
        }
        return selectedRows.length > 1 ? `Export ${selectedRows.length} rows` : `Export`
    }

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
                                    <div className="d-flex gap-1 justify-content-end">
                                        <Select
                                            className="app-search dropdown text-right mw-25"
                                            classNamePrefix="react-select"
                                            placeholder="Filter by status(es)..."
                                            onChange={handleSelectChange}
                                            options={options()}
                                            defaultValue={getOption(parsedQuery.status)}
                                            isMulti
                                            styles={customStyles}
                                            theme={theme => ({
                                                ...theme,
                                                colors: {
                                                    ...theme.colors,
                                                    neutral50: '#98a6ad'
                                                }
                                            })}
                                        />
                                        <Button onClick={() => handleCsvExportClick(data?.vasps)} className="btn btn-light mb-2">{getExportButtonLabel()}</Button>
                                    </div>
                                </Col>
                            </Row>

                            {
                                !isLoading && data ?
                                    <Table
                                        columns={columns}
                                        data={data?.vasps}
                                        pageSize={data?.page_size || 100}
                                        sizePerPageList={sizePerPageList}
                                        isSortable={true}
                                        isSelectable={true}
                                        pagination={true}
                                        isSearchable={true}
                                        theadClass="table-light"
                                        searchBoxClass="mt-2 mb-3"
                                        onSelectedRows={onSelectedRows}
                                    /> : <OvalLoader />
                            }
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
        </React.Fragment>
    );
};

export default VaspsList;
