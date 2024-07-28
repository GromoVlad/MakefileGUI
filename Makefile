LandGenCreateTest_gui:
	echo "<?php namespace App\Console\Commands; use Exception;use Illuminate\Console\Command; class TestRandomClass extends Command{protected signature = 'loc_test'; protected description = 'Command description'; public function __construct(){parent::__construct();} public function handle(){}}" > ~/DATA/Skillbox/courses_generator/app/Console/Commands/TestRandomClass.php
OtherDockerAllDown_gui:
	(echo "[START] Останавливаем Docker-контейнеры") && (cd ~/DATA/Skillbox/courses_generator && make down) && (cd ~/DATA/Skillbox/speakers && make down) && (cd ~/DATA/Skillbox/coral && make down) && (echo "[END] Все Docker-контейнеры остановлены")
LandGenStart_gui:
	(echo "[START] LandGen запускается" && cd ~/DATA/Skillbox/courses_generator && make start && echo "[END] LandGen запущен")
LandGenStop_gui:
	(echo "[START] LandGen останавливается" && cd ~/DATA/Skillbox/courses_generator && make down && echo "[END] LandGen остановлен")
LandGenTerminal_gui:
	gnome-terminal -- docker exec -it php-fpm bash && echo "[END] LandGen терминал открыт"
CoralStart_gui:
	(echo "[START] Coral запускается" && cd ~/DATA/Skillbox/coral && make up && echo "[END] Coral запущен")
CoralStop_gui:
	(echo "[START] Coral останавливается" && cd ~/DATA/Skillbox/coral && make down && echo "[END] Coral остановлен")
CoralTerminal_gui:
	gnome-terminal -- docker exec -it coral-php-fpm-1 bash && echo "[END] Coral терминал открыт"
SpeakerStart_gui:
	(echo "[START] Speaker запускается" && cd ~/DATA/Skillbox/speakers && make up && echo "[END] Speaker запущен")
SpeakerStop_gui:
	(echo "[START] Speaker останавливается" && cd ~/DATA/Skillbox/speakers && make down && echo "[END] Speaker остановлен")
SpeakerTerminal_gui:
	gnome-terminal -- docker exec -it speakers-php-fpm-speakers-1 bash && echo "[END] Speaker терминал открыт"

