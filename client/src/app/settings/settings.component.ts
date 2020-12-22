import { Component, OnInit } from '@angular/core';
import {IEndpointService} from '../services/endpoint.service';
import {ISettings} from '../services/contracts/Settings';

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit {
  public settings: ISettings | undefined;

  constructor(private endpointService: IEndpointService,) { }

  async ngOnInit(): Promise<void> {
    this.settings = await this.endpointService.GetSettings();
  }

  async Save(): Promise<void> {
    if(this.settings){
      await this.endpointService.UpdateSettings(this.settings);
    }
  }
}
