import { useEffect, useState } from "react";
import { copyToClipboard, getBusinessCategory } from "../../utils";
import { Button } from "@chakra-ui/react";
import { Trans, t } from "@lingui/macro";
import { getVaspCategoryValue } from "constants/basic-details";

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
    value: data?.data?.summary?.name
  },
  {
    label: t`Website`,
    value: data?.data?.summary?.website
  },
  {
    label: t`Business Category`,
    value: getBusinessCategory(data?.data?.summary?.business_category)
  },
  {
    label: t`VASP Category`,
    value: getVaspCategoryValue(data?.data?.summary?.vasp_categories)
  },
  {
    label: t`Country of Registration`,
    value: data?.data?.legal_person?.country_of_registration
  },
  {
    label: t`Technical Contact Name`,
    value: data?.data?.contacts?.technical?.name
  },
  {
    label: t`Technical Contact Email`,
    value: data?.data?.contacts?.technical?.email
  },
  {
    label: t`Technical Contact Phone`,
    value: data?.data?.contacts?.technical?.phone
  },
  {
    label: t`Compliance / Legal Contact Name`,
    value: data?.data?.contacts?.legal?.email
  },
  {
    label: t`Compliance / Legal Contact Email`,
    value: data?.data?.contacts?.legal?.email
  },
  {
    label: t`Compliance / Legal Contact Phone`,
    value: data?.data?.contacts?.legal?.phone
  },
  {
    label: t`Administrative Contact Name`,
    value: data?.data?.contacts?.administrative?.name
  },
  {
    label: t`Administrative Contact Email`,
    value: data?.data?.contacts?.administrative?.email
  },
  {
    label: t`Administrative Contact Phone`,
    value: data?.data?.contacts?.administrative?.phone
  },
  {
    label: t`TRISA Endpoint`,
    value: data?.data?.summary?.endpoint
  },
  {
    label: t`Common Name`,
    value: data?.data?.summary?.common_name
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
