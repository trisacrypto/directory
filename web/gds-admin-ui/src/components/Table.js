import React, { useRef, useEffect, forwardRef } from 'react';
import {
    useTable,
    useSortBy,
    usePagination,
    useRowSelect,
    useGlobalFilter,
    useAsyncDebounce,
    useExpanded,
} from 'react-table';
import classNames from 'classnames';
import PropTypes from 'prop-types'
import qs from 'query-string'
import { useHistory } from 'react-router-dom'
import useQuery from 'hooks/useQuery'

// components
import Pagination from './Pagination';

// Define a default UI for filtering
const GlobalFilter = ({ preGlobalFilteredRows, globalFilter, setGlobalFilter, searchBoxClass }) => {
    const count = preGlobalFilteredRows.length;
    const [value, setValue] = React.useState(globalFilter);
    const onChange = useAsyncDebounce((value) => {
        setGlobalFilter(value || undefined);
    }, 200);

    return (
        <div className={classNames(searchBoxClass)}>
            <span className="d-flex align-items-center">
                Search :{' '}
                <input
                    value={value || ''}
                    onChange={(e) => {
                        setValue(e.target.value);
                        onChange(e.target.value);
                    }}
                    placeholder={`${count} records...`}
                    className="form-control w-auto ms-1"
                />
            </span>
        </div>
    );
};

const IndeterminateCheckbox = forwardRef(({ indeterminate, ...rest }, ref) => {
    const defaultRef = useRef();
    const resolvedRef = ref || defaultRef;

    useEffect(() => {
        resolvedRef.current.indeterminate = indeterminate;
    }, [resolvedRef, indeterminate]);

    return (
        <>
            <div className="form-check">
                <input type="checkbox" className="form-check-input" ref={resolvedRef} {...rest} />
                <label htmlFor="form-check-input" className="form-check-label"></label>
            </div>
        </>
    );
});


const Table = (props) => {
    const isSearchable = props['isSearchable'] || false;
    const isSortable = props['isSortable'] || false;
    const pagination = props['pagination'] || false;
    const isSelectable = props['isSelectable'] || false;
    const isExpandable = props['isExpandable'] || false;
    const hiddenColumns = props['hiddenColumns'] || []
    const query = useQuery()
    const history = useHistory()

    const dataTable = useTable(
        {
            columns: props['columns'],
            data: props['data'],
            initialState: {
                pageSize: props['pageSize'] || 10, hiddenColumns,
                sortBy: [{
                    id: query.get('sortBy'),
                    desc: JSON.parse(query.get('desc'))
                }]
            },
            autoResetSortBy: false

        },
        isSearchable && useGlobalFilter,
        isSortable && useSortBy,
        isExpandable && useExpanded,
        pagination && usePagination,
        isSelectable && useRowSelect,
        (hooks) => {
            isSelectable &&
                hooks.visibleColumns.push((columns) => [
                    {
                        id: 'selection',
                        Header: ({ getToggleAllPageRowsSelectedProps }) => (
                            <div>
                                <IndeterminateCheckbox {...getToggleAllPageRowsSelectedProps()} />
                            </div>
                        ),
                        Cell: ({ row }) => (
                            <div>
                                <IndeterminateCheckbox {...row.getToggleRowSelectedProps()} />
                            </div>
                        ),
                    },
                    ...columns,
                ]);

            isExpandable &&
                hooks.visibleColumns.push((columns) => [
                    {
                        id: 'expander',
                        Header: ({ getToggleAllRowsExpandedProps, isAllRowsExpanded }) => (
                            <span {...getToggleAllRowsExpandedProps()}>{isAllRowsExpanded ? '-' : '+'}</span>
                        ),
                        Cell: ({ row }) =>
                            row.canExpand ? (
                                <span
                                    {...row.getToggleRowExpandedProps({
                                        style: {
                                            paddingLeft: `${row.depth * 2}rem`,
                                        },
                                    })}>
                                    {row.isExpanded ? '-' : '+'}
                                </span>
                            ) : null,
                    },
                    ...columns,
                ]);
        }
    );

    const sortBy = dataTable?.state.sortBy[0]

    React.useEffect(() => {
        props.onSelectedRows(dataTable.selectedFlatRows);

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [dataTable.selectedFlatRows])

    React.useEffect(() => {
        if (sortBy?.id) {
            const query = { sortBy: sortBy?.id || 'name', desc: sortBy?.desc || false }
            history.push({ search: qs.stringify(query) })
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [sortBy?.id, sortBy?.desc])

    React.useEffect(() => {

    })

    const rows = pagination ? dataTable.page : dataTable.rows;

    return (
        <>
            {isSearchable && (
                <GlobalFilter
                    preGlobalFilteredRows={dataTable.preGlobalFilteredRows}
                    globalFilter={dataTable.state.globalFilter}
                    setGlobalFilter={dataTable.setGlobalFilter}
                    searchBoxClass={props['searchBoxClass']}
                />
            )}

            <div className="table-responsive">
                <table
                    {...dataTable.getTableProps()}
                    className={classNames('table table-centered react-table', props['tableClass'])}>
                    <thead className={props['theadClass']}>
                        {dataTable.headerGroups.map((headerGroup) => (
                            <tr {...headerGroup.getHeaderGroupProps()}>
                                {headerGroup.headers.map((column) => (
                                    <th
                                        {...column.getHeaderProps(column.sort && column.getSortByToggleProps())}
                                        className={classNames({
                                            sorting_desc: column.isSortedDesc === true,
                                            sorting_asc: column.isSortedDesc === false,
                                            sortable: column.sort === true,
                                        })}>
                                        {column.render('Header')}
                                    </th>
                                ))}
                            </tr>
                        ))}
                    </thead>

                    <tbody {...dataTable.getTableBodyProps()}>
                        {
                            !rows.length ? <tr><td>No vasps available</td></tr> : null
                        }
                        {(rows || []).map((row, i) => {
                            dataTable.prepareRow(row);
                            return (
                                <tr {...row.getRowProps()}>
                                    {row.cells.map((cell) => {
                                        return <td {...cell.getCellProps()}>{cell.render('Cell')}</td>;
                                    })}
                                </tr>
                            );
                        })}
                    </tbody>
                </table>
            </div>

            {pagination && <Pagination tableProps={dataTable} sizePerPageList={props['sizePerPageList']} />}
        </>
    );
};

Table.propTypes = {
    isSearchable: PropTypes.bool,
    isSortable: PropTypes.bool,
    pagination: PropTypes.bool,
    isSelectable: PropTypes.bool,
    isExpandable: PropTypes.bool,
    hiddenColumns: PropTypes.arrayOf(PropTypes.string),
    columns: PropTypes.arrayOf(PropTypes.object),
    pageSize: PropTypes.number,
    data: PropTypes.any,
    sizePerPageList: PropTypes.arrayOf(PropTypes.shape({
        text: PropTypes.string.isRequired,
        value: PropTypes.number.isRequired
    })),
    onSelectedRows: PropTypes.func,
    theadClass: PropTypes.string,
    tableClass: PropTypes.string,
    searchBoxClass: PropTypes.string
}

export default Table;
