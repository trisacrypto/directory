import React from 'react';
import { Stack, Box, Text, Heading, Flex, chakra } from '@chakra-ui/react';
import { t } from '@lingui/macro';

interface StatCardProps {
  title: string;
  number: number;
}
const StatCard = ({ title, number }: StatCardProps) => {
  return (
    <Box
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      color={'#252733'}
      textAlign={'center'}
      fontSize={18}
      p={5}
      mt={10}
      px={5}>
      <Stack>
        <Text as="h1" textAlign={'center'} fontSize={20} fontWeight={'bold'}>
          {title}
        </Text>
        <Text fontSize={40} pt={3} fontWeight={'bold'}>
          {number}
        </Text>
      </Stack>
    </Box>
  );
};
StatCard.defaultProps = {
  title: t`Verified VASPs`,
  number: 0
};

export default StatCard;
