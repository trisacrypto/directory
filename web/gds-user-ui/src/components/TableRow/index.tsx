import RowItem from './RowItem';
import { Td } from '@chakra-ui/react';
import { ReactNode } from 'react';

type TableRowProps = {
  row: { [k: string]: ReactNode };
};

function TableRow({ row }: TableRowProps) {
  return (
    <RowItem>
      {Object.entries(row).map(([k, v]) => (
        <Td key={k}>{v}</Td>
      ))}
    </RowItem>
  );
}

export default TableRow;
