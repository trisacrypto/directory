import { HStack, Flex, IconButton, Box, Text } from '@chakra-ui/react';
import React from 'react';
import { FiDownload } from 'react-icons/fi';

function formatBytes(bytes: number, decimals = 2) {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / k ** i).toFixed(dm)) + ' ' + sizes[i];
}

function getBase64Size(str: string) {
  const buffer = Buffer.from(`${str}`, 'base64');
  return buffer.length;
}

type FileCardProps = {
  name: string;
  file?: string;
  ext: string;
  onDownload: () => void;
};

export function FileCard({ name, file, ext, onDownload }: FileCardProps) {
  const fileSize = React.useMemo(() => (file ? formatBytes(getBase64Size(file)) : ''), [file]);

  return (
    <HStack
      border="1px solid #00000094"
      borderRadius="10px"
      p={3}
      alignItems="center!important"
      spacing={5}>
      <Flex gap={2}>
        <Flex
          bg="#23A7E04D"
          borderRadius="10px"
          fontWeight={700}
          color="rgba(85, 81, 81, 0.83)"
          justifyContent="center"
          alignItems="center"
          p={2}>
          {ext}
        </Flex>
        <Box>
          <Text fontWeight={700}>{name}</Text>
          <Text color="gray.600" fontSize="sm">
            {fileSize || '0 Kb'}
          </Text>
        </Box>
      </Flex>
      <IconButton
        variant="ghost"
        fontSize="30px"
        color="blue"
        onClick={onDownload}
        p={3}
        icon={<FiDownload />}
        aria-label="download"
        disabled={!file}
      />
    </HStack>
  );
}
