import {DateAdapter} from "@angular/material/core";

export class IPAddress {
  IP: string;
  Zone: string;
}
export class Lease {
  macAddress: string;
  expireTime: Date;
  ipAddress: IPAddress;
  name: string;
  clientId: string;
}
