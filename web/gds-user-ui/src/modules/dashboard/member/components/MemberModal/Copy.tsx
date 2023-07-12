import { useEffect, useState } from "react";
import { copyToClipboard } from "../../utils";
import { Button } from "@chakra-ui/react";
import { Trans, t } from "@lingui/macro";
import { BUSINESS_CATEGORY, getVaspCategoryValue } from "constants/basic-details";

interface MemberDetailTableHeader {
  label: string;
  value: string;
}

type CopyProps = {
    data: any;
};

function Copy({ data }: CopyProps) {
const memberDetailTableHeader = [
  {
    label: t`Name`,
    value: data?.summary?.name || 'N/A'
  },
  {
    label: t`Website`,
    value: data?.summary?.website || 'N/A'
  },
  {
    label: t`Business Category`,
    value: BUSINESS_CATEGORY[data?.summary?.business_category as keyof typeof BUSINESS_CATEGORY] || 'N/A'
  },
  {
    label: t`VASP Category`,
    value: getVaspCategoryValue(data?.summary?.vasp_categories) || 'N/A'
  },
  {
    label: t`Country of Registration`,
    value: data?.legal_person?.country_of_registration || 'N/A'
  },
  {
    label: t`Technical Contact Name`,
    value: data?.contacts?.technical?.name || 'N/A'
  },
  {
    label: t`Technical Contact Email`,
    value: data?.contacts?.technical?.email || 'N/A'
  },
  {
    label: t`Technical Contact Phone`,
    value: data?.contacts?.technical?.phone || 'N/A'
  },
  {
    label: t`Compliance / Legal Contact Name`,
    value: data?.contacts?.legal?.email || 'N/A'
  },
  {
    label: t`Compliance / Legal Contact Email`,
    value: data?.contacts?.legal?.email || 'N/A'
  },
  {
    label: t`Compliance / Legal Contact Phone`,
    value: data?.contacts?.legal?.phone || 'N/A'
  },
  {
    label: t`Administrative Contact Name`,
    value: data?.contacts?.administrative?.name || 'N/A'
  },
  {
    label: t`Administrative Contact Email`,
    value: data?.contacts?.administrative?.email || 'N/A'
  },
  {
    label: t`Administrative Contact Phone`,
    value: data?.contacts?.administrative?.phone || 'N/A'
  },
  {
    label: t`TRISA Endpoint`,
    value: data?.summary?.endpoint || 'N/A'
  },
  {
    label: t`Common Name`,
    value: data?.summary?.common_name || 'N/A'
  },
];
    const [copied, setCopied] = useState(false);

    useEffect(() => {
        if (copied) {
            const timeout = setTimeout(() => {
                setCopied(false);
            }, 2000);
            return () => clearTimeout(timeout);
        }
    }, [copied]);


    const handleCopy = async () => {
        await copyToClipboard(memberDetailTableHeader as MemberDetailTableHeader[]);
        setCopied(true);
    };
    return copied ? (
        <Button bg={'#FF7A59'} color={'white'}>
            <Trans>Copied</Trans>
        </Button>
    ) : (
        <Button bg={'#FF7A59'} color={'white'} onClick={handleCopy}>
            <Trans>Copy</Trans>
        </Button>
    );
}

export default Copy;
