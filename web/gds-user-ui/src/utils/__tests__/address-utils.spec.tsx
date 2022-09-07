import {
  hasAddressField,
  hasAddressFieldAndLine,
  hasAddressLine,
  isValidIvmsAddress,
  renderAddress
} from 'utils/address-utils';

describe('Address Utils', () => {
  describe('isValidIvmsAddress', () => {
    it('should be a valid ivms address', () => {
      const ivms101Address = { country: 'US', address_type: 'GEOG' };
      expect(isValidIvmsAddress(ivms101Address)).toBe(true);
    });

    it('should not be a valid ivms address if there are no country', () => {
      const ivms101Address = { address_type: 'GEOG' };
      expect(isValidIvmsAddress(ivms101Address)).toBe(false);
    });

    it('should not be a valid ivms address if there are no address type', () => {
      const ivms101Address = { country: 'US' };
      expect(isValidIvmsAddress(ivms101Address)).toBe(false);
    });
  });

  describe('hasAddressLine', () => {
    it('should return true when there are address line and is a valid IVMS address', () => {
      const ivms101Address = {
        country: 'US',
        address_type: 'GEOG',
        address_line: ['Gangnam-gu, Gangbuck-do']
      };

      expect(hasAddressLine(ivms101Address)).toBe(true);
    });

    it('should return false when there are address line and is not a valid IVMS address', () => {
      const ivms101Address = {
        country: 'US',
        address_type: '',
        address_line: ['Gangnam-gu, Gangbuck-do']
      };

      expect(hasAddressLine(ivms101Address)).toBe(false);
    });

    it('should return false when there are no address line and is a valid IVMS address', () => {
      const ivms101Address = {
        country: 'US',
        address_type: 'GEOG',
        address_line: []
      };

      expect(hasAddressLine(ivms101Address)).toBe(false);
    });
  });

  describe('hasAddressField', () => {
    it('should have address field', () => {
      const ivms101Address = {
        country: 'US',
        address_type: 'GEOG',
        street_name: 'Schroeder Isle',
        building_number: '5786'
      };

      expect(hasAddressField(ivms101Address)).toBe(true);
    });

    it('should have address field', () => {
      const ivms101Address = {
        country: 'US',
        address_type: 'GEOG',
        street_name: 'Schroeder Isle',
        building_name: 'Little Summit'
      };

      expect(hasAddressField(ivms101Address)).toBe(true);
    });

    it('should return false when we do not have building name or building number', () => {
      const ivms101Address = {
        country: 'US',
        address_type: 'GEOG',
        street_name: 'Schroeder Isle',
        building_name: ''
      };

      expect(hasAddressField(ivms101Address)).toBe(false);
    });

    it('should return false when we do not have building name or building number', () => {
      const ivms101Address = {
        country: 'US',
        address_type: 'GEOG',
        street_name: 'Schroeder Isle',
        building_number: ''
      };

      expect(hasAddressField(ivms101Address)).toBe(false);
    });

    it('should return false when we do not have building name or building number', () => {
      const ivms101Address = {
        country: 'US',
        address_type: 'GEOG',
        street_name: 'Schroeder Isle',
        building_number: 'Little Summit',
        address_line: ['Gangnam-gu, Gangbuck-do']
      };

      expect(hasAddressField(ivms101Address)).toBe(false);
    });
  });

  describe('renderAddress', () => {
    it('should be unparseable address', () => {
      const ivms101Address = {};

      expect(renderAddress(ivms101Address)).toEqual(<div>Unparseable Address</div>);
    });

    it('should match inline snapshot', () => {
      const ivms101Address = {
        country: 'US',
        address_type: 'GEOG',
        street_name: 'Schroeder Isle',
        building_name: 'Little Summit',
        address_line: ['Gangnam-gu, Gangbuck-do']
      };

      expect(renderAddress(ivms101Address)).toMatchInlineSnapshot(`
        <div
          data-testid="addressLine"
        >
          <div>
            Gangnam-gu, Gangbuck-do
             
          </div>
          <div>
            undefined, undefined undefined
          </div>
          <div>
            US
          </div>
        </div>
      `);
    });

    it('should match inline snapshot', () => {
      const ivms101Address = {
        country: 'US',
        address_type: 'GEOG',
        street_name: 'Schroeder Isle',
        building_name: 'Little Summit'
      };

      expect(renderAddress(ivms101Address)).toMatchInlineSnapshot(`
        <div
          data-testid="addressField"
        >
           
          Schroeder Isle
          <br />
          <React.Fragment>
             
             
            Little Summit
             
            <br />
          </React.Fragment>
          US
        </div>
      `);
    });
  });
});
