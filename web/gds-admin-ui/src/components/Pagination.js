// @flow
import React from 'react';

type PaginationProps = {
    tableProps: any,
    sizePerPageList: {
        text: string,
        value: number,
    }[],
};

const Pagination = ({ tableProps, sizePerPageList }: PaginationProps): React$Element<React$FragmentType> => {
    return (
        <>
            <div className="d-lg-flex align-items-center text-center pb-1">
                <label className="me-1">Display :</label>
                <select
                    value={tableProps.state.pageSize}
                    onChange={(e) => {
                        tableProps.setPageSize(Number(e.target.value));
                    }}
                    className="form-select d-inline-block w-auto">
                    {sizePerPageList.map((pageSize, index) => {
                        return (
                            <option key={index} value={pageSize.value}>
                                {pageSize.text}
                            </option>
                        );
                    })}
                </select>

                <span className="mx-3">
                    Page{' '}
                    <strong>
                        {tableProps.state.pageIndex + 1} of {tableProps.pageOptions.length}
                    </strong>{' '}
                </span>

                <span className="d-inline-block align-items-center text-sm-start text-center my-sm-0 my-2">
                    <label>Go to page : </label>
                    <input
                        type="number"
                        defaultValue={tableProps.state.pageIndex + 1}
                        onChange={(e) => {
                            const page = e.target.value ? Number(e.target.value) - 1 : 0;
                            tableProps.gotoPage(page);
                        }}
                        className="form-control w-25 ms-1 d-inline-block"
                    />
                </span>

                <div className="pagination pagination-rounded d-inline-flex ms-auto">
                    <button
                        onClick={() => tableProps.gotoPage(0)}
                        className="page-link"
                        disabled={!tableProps.canPreviousPage}>
                        {'<<'}
                    </button>{' '}
                    <button
                        onClick={() => tableProps.previousPage()}
                        className="page-link"
                        disabled={!tableProps.canPreviousPage}>
                        {'<'}
                    </button>{' '}
                    <button
                        onClick={() => tableProps.nextPage()}
                        className="page-link"
                        disabled={!tableProps.canNextPage}>
                        {'>'}
                    </button>{' '}
                    <button
                        onClick={() => tableProps.gotoPage(tableProps.pageCount - 1)}
                        className="page-link"
                        disabled={!tableProps.canNextPage}>
                        {'>>'}
                    </button>{' '}
                </div>
            </div>
        </>
    );
};

export default Pagination;
