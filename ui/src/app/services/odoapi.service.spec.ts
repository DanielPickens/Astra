import { TestBed } from '@angular/core/testing';

import { astraapiService } from './astraapi.service';

describe('astraapiService', () => {
  let service: astraapiService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(astraapiService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
