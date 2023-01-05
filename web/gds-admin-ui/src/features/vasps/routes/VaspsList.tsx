import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import qs from 'query-string';
import React from 'react';
import { Button, Card, Col, Row } from 'react-bootstrap';
import { useDispatch, useSelector } from 'react-redux';
import { useHistory, useLocation } from 'react-router-dom';
import Select from 'react-select';

import OvalLoader from '@/components/OvalLoader';
import PageTitle from '@/components/PageTitle';
import Table from '@/components/Table';
import useSafeDispatch from '@/hooks/useSafeDispatch';
import { fetchVasps } from '@/redux/dashboard/actions';
import { getAllVasps, getVaspsLoadingState } from '@/redux/selectors/vasps';
import { exportToCsv } from '@/utils';

import CertificateExpirationColumn from '../components/lists/CertificateExpirationColumn';
import LastUpdatedColumn from '../components/lists/LastUpdatedColumn';
import NameColumn from '../components/lists/NameColumn';
import StatusColumn from '../components/lists/StatusColumn';
import { StatusLabel } from '@/constants';

dayjs.extend(relativeTime);

const columns = [
    {
        Header: 'Name',
        accessor: 'name',
        sort: true,
        Cell: NameColumn,
    },
    {
        Header: 'Common Name',
        accessor: 'common_name',
        sort: true,
        Cell: <></>,
    },
    {
        Header: 'Status',
        accessor: 'verification_status',
        sort: true,
        Cell: StatusColumn,
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
        Cell: CertificateExpirationColumn,
    },
];

const options = () => Object.entries(StatusLabel).map(([k, v]) => ({ value: k, label: v }));

const getOption = (key: any) => {
    const _options = options();
    let opt: any[] = [];

    for (const option of _options) {
        if (typeof key === 'string' && option.value.toLowerCase() === key) {
            opt = [option];
        }

        if (Array.isArray(key)) {
            for (const k of key) {
                if (k === option.value.toLowerCase()) {
                    opt.push(option);
                }
            }
        }
    }
    return opt;
};

const customStyles = {
    control: (styles: any) => ({ ...styles, paddingLeft: '9px !important' }),
};

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
        value: 100,
    },
];

const VaspsList = () => {
    const location = useLocation();
    const parsedQuery = qs.parse(location.search);
    const query = qs.stringify(parsedQuery);

    const [queryParams, setQueryParams] = React.useState(query);
    const dispatch = useDispatch();
    const safeDispatch = useSafeDispatch(dispatch);
    const data = useSelector(getAllVasps);
    const isLoading = useSelector(getVaspsLoadingState);
    const history = useHistory();
    const [selectedRows, setSelectedData] = React.useState<any>([]);

    React.useEffect(() => {
        safeDispatch(fetchVasps({ queryParams }));
    }, [safeDispatch, queryParams]);

    const getQueryString = (arr: any) => {
        const queries = arr && Array.isArray(arr) ? arr.map((v) => v.value.toLowerCase()) : '';
        return qs.stringify({ status: queries });
    };

    const handleSelectChange = (value: any) => {
        const params = getQueryString(value);
        setQueryParams(params);
        history.push({
            pathname: '/vasps',
            search: params,
        });
    };

    const handleCsvExportClick = (rows: any) => (selectedRows.length ? exportToCsv(selectedRows) : exportToCsv(rows));

    const onSelectedRows = (rows: any) => {
        const mappedRows = rows.map((r: any) => r.original);
        setSelectedData([...mappedRows]);
    };

    const getExportButtonLabel = () => {
        if (selectedRows.length === 1) {
            return 'Export 1 row';
        }
        return selectedRows.length > 1 ? `Export ${selectedRows.length} rows` : 'Export';
    };

    return (
        <>
            <PageTitle
                breadCrumbItems={[{ label: 'List', path: '/vasps', active: true }]}
                title="All Registered VASPs"
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
                                            theme={(theme: any) => ({
                                                ...theme,
                                                colors: {
                                                    ...theme.colors,
                                                    neutral50: '#98a6ad',
                                                },
                                            })}
                                        />
                                        <Button
                                            onClick={() => handleCsvExportClick(data?.vasps)}
                                            className="btn btn-light mb-2">
                                            {getExportButtonLabel()}
                                        </Button>
                                    </div>
                                </Col>
                            </Row>

                            {!isLoading && data ? (
                                <Table
                                    columns={columns}
                                    data={data?.vasps}
                                    pageSize={data?.page_size || 100}
                                    sizePerPageList={sizePerPageList}
                                    hiddenColumns={['common_name']}
                                    isSortable
                                    isSelectable
                                    pagination
                                    isSearchable
                                    theadClass="table-light"
                                    searchBoxClass="mt-2 mb-3"
                                    onSelectedRows={onSelectedRows}
                                />
                            ) : (
                                <OvalLoader title={''} />
                            )}
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
        </>
    );
};

export default VaspsList;
