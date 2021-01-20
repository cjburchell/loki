import { Component, OnInit } from '@angular/core';
import {IEndpointService} from '../services/endpoint.service';
import {ISettings} from '../services/contracts/Settings';
import {__await} from 'tslib';
import {Router} from '@angular/router';

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit {
  public settings: ISettings | undefined;

  public errorMsg = '';

  constructor(private endpointService: IEndpointService,
              private router: Router) {
  }

  async ngOnInit(): Promise<void> {
    try{
      this.settings = await this.endpointService.GetSettings();
    }catch (e) {
      this.errorMsg = `Error loading settings: ${e.status} (${e.statusText})`;
    }
  }

  public async Save(): Promise<void> {
    if (this.settings) {
      try {
        await this.endpointService.UpdateSettings(this.settings);
        this.router.navigate([`/endpoints`]).then(() => {});
      } catch (e) {
        this.errorMsg = `Unable to save settings: ${e.status} (${e.statusText})`;
      }
    }
  }
}
