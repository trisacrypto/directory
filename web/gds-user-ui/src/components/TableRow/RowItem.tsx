import React from 'react';
import { Tr } from '@chakra-ui/react';

const RowItem: React.FC<{ children: React.ReactNode }> = ({ children }) => <Tr>{children}</Tr>;

export default RowItem;
