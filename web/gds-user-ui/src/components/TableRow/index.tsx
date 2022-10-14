import RowItem from './RowItem';
import { Td } from '@chakra-ui/react';
import { ReactNode } from 'react';

type TableRowProps<T> = {
  row: T | { [k: string]: ReactNode };
};

function TableRow<T>({ row }: TableRowProps<T>) {
  return (
    <RowItem>
      {Object.entries(row).map(([k, v]) => (
        <Td key={k}>{v}</Td>
      ))}
    </RowItem>
  );
}

export default TableRow;
